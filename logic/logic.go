package logic

import (
	db2 "abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/db/postgres"
	"abbysoft/gardarike-online/model"
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
}

type Config struct {
	AFKTimeout           time.Duration
	ChatMessageMaxLength int
}

func NewLogic(generator TerrainGenerator, eventsChan chan model.EventWrapper, dbConfig postgres.Config, config Config) (*SimpleLogic, error) {
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
	}

	logic.resourceManager = NewResourceManager(logic)

	logic.log.Info("Initialize logic...")
	if err := logic.init(generator); err != nil {
		return nil, fmt.Errorf("failed to init data from the DB: %w", err)
	}
	logic.log.Info("Logic initialization is done")

	logic.log.Info("Running game loop")
	go logic.gameLoop()

	return logic, nil
}

func (s *SimpleLogic) SaveGameMap() error {
	return s.db.SaveOrUpdate(model.NewWorldMapChunkFromRPC(s.GameMap), true)
}

func (s *SimpleLogic) generateGameMap(generator TerrainGenerator) error {
	s.log.Info("Map not found, generating it...")

	terrain := generator.GenerateTerrain(mapChunkSize, mapChunkSize)
	s.GameMap = rpc.WorldMapChunk{
		X:       1,
		Y:       1,
		Width:   model.MapChunkSize,
		Height:  model.MapChunkSize,
		Data:    terrain,
		Towns:   []*rpc.Town{},
		Trees:   0,
		Stones:  0,
		Animals: 0,
		Plants:  0,
	}

	if err := s.SaveGameMap(); err != nil {
		return fmt.Errorf("failed to save game map: %w", err)
	}

	return nil
}

func (s *SimpleLogic) loadOrGenerateGameMap(generator TerrainGenerator) error {
	mapChunk, err := s.db.GetMapChunk(0, 0)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to load stored map: %w", err)
	} else if err != nil {
		return s.generateGameMap(generator)
	}

	rpcChunk, err := mapChunk.ToRPC()
	if err != nil {
		return fmt.Errorf("failed to convert map chunk to rpc struct: %w", err)
	}

	s.GameMap = *rpcChunk
	return nil
}

func (s *SimpleLogic) init(generator TerrainGenerator) error {
	s.log.Info("Initializing game map")
	if err := s.loadOrGenerateGameMap(generator); err != nil {
		return fmt.Errorf("failed to init game map: %w", err)
	}

	s.log.
		WithField("points", len(s.GameMap.Data)).
		WithField("trees", s.GameMap.Trees).
		WithField("stones", s.GameMap.Stones).
		WithField("animas", s.GameMap.Animals).
		WithField("plants", s.GameMap.Plants).
		Info("Game map initialized")

	s.log.Info("Loading town locations...")

	// Load town locations
	// TODO

	s.log.Infof("Loaded %d towns on the map", 0)
	return nil
}

func (s *SimpleLogic) GetWorldMap(_ *PlayerSession, request *rpc.GetWorldMapRequest) (*rpc.GetWorldMapResponse, model.Error) {
	s.log.WithField("location", request.GetLocation()).
		WithField("sessionID", request.GetSessionID()).
		Infof("GetMap request")

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
		return nil, model.ErrNotAuthorized
	}

	session.SelectedCharacter = &char
	s.log.WithFields(logrus.Fields{
		"sessionID": request.GetSessionID(),
		"character": char,
	}).Info("User selected character")

	return &rpc.SelectCharacterResponse{}, nil
}
