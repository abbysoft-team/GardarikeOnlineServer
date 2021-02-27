package logic

import (
	db2 "abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/db/postgres"
	"abbysoft/gardarike-online/generation"
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	mapChunkSize = 100
)

type Logic interface {
	GetWorldMap(session *PlayerSession, request *rpc.GetWorldMapRequest) (*rpc.GetWorldMapResponse, model.Error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, model.Error)
	SelectCharacter(session *PlayerSession, request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error)
	SendChatMessage(session *PlayerSession, request *rpc.SendChatMessageRequest) (*rpc.SendChatMessageResponse, model.Error)
	GetChatHistory(session *PlayerSession, request *rpc.GetChatHistoryRequest) (*rpc.GetChatHistoryResponse, model.Error)
	GetWorkDistribution(session *PlayerSession, request *rpc.GetWorkDistributionRequest) (*rpc.GetWorkDistributionResponse, model.Error)
	CreateAccount(session *PlayerSession, request *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, model.Error)
	CreateCharacter(session *PlayerSession, request *rpc.CreateCharacterRequest) (*rpc.CreateCharacterResponse, model.Error)
	GetResources(session *PlayerSession, request *rpc.GetResourcesRequest) (*rpc.GetResourcesResponse, model.Error)
	PlaceTown(session *PlayerSession, request *rpc.PlaceTownRequest) (*rpc.PlaceTownResponse, model.Error)
}

type SimpleLogic struct {
	GameMap         rpc.WorldMapChunk
	GameMapMutex    sync.Mutex
	db              db2.Database
	log             *logrus.Entry
	sessions        map[string]*PlayerSession
	EventsChan      chan model.EventWrapper
	config          Config
	resourceManager ResourceManager
	generator       generation.TerrainGenerator
}

type Config struct {
	AFKTimeout           time.Duration
	ChatMessageMaxLength int
	WaterLevel           float32
	ChunkSize            int
	AlwaysRegenerateMap  bool
}

func NewLogic(generator generation.TerrainGenerator, eventsChan chan model.EventWrapper, dbConfig postgres.Config, config Config) (*SimpleLogic, error) {
	database, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	logic := &SimpleLogic{
		db:         database,
		log:        logrus.WithField("module", "logic"),
		sessions:   make(map[string]*PlayerSession),
		EventsChan: eventsChan,
		config:     config,
		generator:  generator,
	}

	logic.resourceManager = NewResourceManager(logic)

	logic.log.WithField("config", config).Info("Initialize logic...")

	if err := logic.init(); err != nil {
		return nil, fmt.Errorf("failed to init data from the DB: %w", err)
	}
	logic.log.Info("Logic initialization is done")

	logic.log.Info("Running game loop")
	go logic.gameLoop()

	return logic, nil
}

func (s *SimpleLogic) SaveGameMap() error {
	chunk, err := model.NewWorldMapChunkFromRPC(s.GameMap)
	if err != nil {
		return err
	}

	return s.db.SaveOrUpdate(chunk, true)
}

func (s *SimpleLogic) generateGameMap() error {
	s.log.Info("Map not found, generating it...")

	terrain := s.generator.GenerateTerrain(s.config.ChunkSize, s.config.ChunkSize, 0, 0)
	s.GameMap = rpc.WorldMapChunk{
		X:          0,
		Y:          0,
		Width:      int32(s.config.ChunkSize),
		Height:     int32(s.config.ChunkSize),
		Data:       terrain,
		Towns:      []*rpc.Town{},
		Trees:      0,
		Stones:     0,
		Animals:    0,
		Plants:     0,
		WaterLevel: s.config.WaterLevel,
	}

	if err := s.SaveGameMap(); err != nil {
		return fmt.Errorf("failed to save game map: %w", err)
	}

	return nil
}

func (s *SimpleLogic) loadOrGenerateGameMap() error {
	mapChunk, err := s.db.GetMapChunk(0, 0)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to load stored map: %w", err)
	} else if err != nil {
		return s.generateGameMap()
	}

	rpcChunk, err := mapChunk.ToRPC()
	if err != nil {
		return fmt.Errorf("failed to convert map chunk to rpc struct: %w", err)
	}

	rpcChunk.WaterLevel = s.config.WaterLevel
	s.GameMap = *rpcChunk
	return nil
}

func (s *SimpleLogic) init() error {
	s.log.Info("Initializing game map")
	if err := s.loadOrGenerateGameMap(); err != nil {
		return fmt.Errorf("failed to init game map: %w", err)
	}

	s.log.
		WithField("points", len(s.GameMap.Data)).
		WithField("trees", s.GameMap.Trees).
		WithField("stones", s.GameMap.Stones).
		WithField("animas", s.GameMap.Animals).
		WithField("plants", s.GameMap.Plants).
		Info("Game map initialized")

	s.log.Info("Loading towns...")

	towns, err := s.db.GetAllTowns()
	if err != nil {
		return fmt.Errorf("failed to load towns: %w", err)
	}

	for _, town := range towns {
		s.GameMap.Towns = append(s.GameMap.Towns, town.ToRPC())
	}

	s.log.Infof("Loaded %d towns on the map", len(s.GameMap.Towns))

	return nil
}

func (s *SimpleLogic) GetWorldMap(_ *PlayerSession, request *rpc.GetWorldMapRequest) (*rpc.GetWorldMapResponse, model.Error) {
	s.log.WithField("location", request.GetLocation()).
		WithField("sessionID", request.GetSessionID()).
		Infof("GetMap request")

	if s.config.AlwaysRegenerateMap {
		s.generator.SetSeed(time.Now().UnixNano())

		if err := s.generateGameMap(); err != nil {
			s.log.WithError(err).Error("Failed to regenerate game map")
			return nil, model.ErrInternalServerError
		}
	}

	return &rpc.GetWorldMapResponse{Map: &s.GameMap}, nil
}

func (s *SimpleLogic) SelectCharacter(session *PlayerSession, request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error) {
	s.log.WithField("characterID", request.GetCharacterID()).
		WithField("sessionID", request.GetSessionID()).
		Infof("SelectCharacter request")

	char, err := s.db.GetCharacter(request.GetCharacterID())
	if err != nil {
		return nil, model.ErrCharacterNotFound
	}

	if char.AccountID != session.AccountID {
		return nil, model.ErrForbidden
	}

	towns, err := s.db.GetTowns(char.Name)
	if err != nil {
		s.log.WithError(err).Error("Failed to get character's towns")
	} else {
		char.Towns = towns
	}

	resources, err := s.db.GetResources(char.ID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get character's resources")
	} else {
		char.Resources = resources
	}

	session.SelectedCharacter = &char
	s.log.WithFields(logrus.Fields{
		"sessionID": request.GetSessionID(),
		"character": char,
	}).Info("User selected character")

	s.EventsChan <- model.NewSystemChatMessageEvent(consts.MessageCharacterAuthorized(char.Name))

	response := &rpc.SelectCharacterResponse{}
	for _, town := range char.Towns {
		response.Towns = append(response.Towns, town.ToRPC())
	}

	return response, nil
}
