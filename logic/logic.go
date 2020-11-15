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
	GetMap(session *PlayerSession, request *rpc.GetMapRequest) (*rpc.GetMapResponse, model.Error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, model.Error)
	SelectCharacter(session *PlayerSession, request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error)
	PlaceBuilding(session *PlayerSession, request *rpc.PlaceBuildingRequest) (*rpc.PlaceBuildingResponse, model.Error)
	SendChatMessage(session *PlayerSession, request *rpc.SendChatMessageRequest) (*rpc.SendChatMessageResponse, model.Error)
	GetChatHistory(session *PlayerSession, request *rpc.GetChatHistoryRequest) (*rpc.GetChatHistoryResponse, model.Error)
}

type SimpleLogic struct {
	GameMap         rpc.Map
	GameMapMutex    sync.Mutex
	buildings       map[int]model.Building
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
		buildings:  make(map[int]model.Building),
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
	modelMap, err := model.NewMapChunkFrom(s.GameMap)
	if err != nil {
		return fmt.Errorf("failed to create new map chunk: %w", err)
	}

	return s.db.SaveOrUpdate(modelMap)
}

func (s *SimpleLogic) generateGameMap(generator TerrainGenerator) error {
	s.log.Info("Map not found, generating it...")

	terrain := generator.GenerateTerrain(mapChunkSize, mapChunkSize)
	s.GameMap = rpc.Map{
		Width:      model.MapChunkSize,
		Height:     model.MapChunkSize,
		Points:     terrain,
		Buildings:  nil,
		TreesCount: 0,
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

	gameMap, err := mapChunk.ToRPC()
	if err != nil {
		return fmt.Errorf("failed to serialize map chunk: %w", err)
	}

	s.GameMap = *gameMap
	return nil
}

func (s *SimpleLogic) init(generator TerrainGenerator) error {
	s.log.Info("Initializing game map")
	if err := s.loadOrGenerateGameMap(generator); err != nil {
		return fmt.Errorf("failed to init game map: %w", err)
	}

	s.log.
		WithField("points", len(s.GameMap.Points)).
		WithField("treesCount", s.GameMap.TreesCount).
		Info("Game map initialized")

	s.log.Info("Loading building locations...")

	// Load building locations
	buildingLocations, err := s.db.GetBuildingLocations()
	if err != nil {
		return fmt.Errorf("failed to init building locations: %w", err)
	}

	var rpcBuildings []*rpc.Building
	for _, building := range buildingLocations {
		rpcBuildings = append(rpcBuildings, building.ToRPC())
	}

	s.GameMap.Buildings = rpcBuildings

	s.log.Infof("Loaded %d buildings on the map", len(s.GameMap.Buildings))
	s.log.Infof("Loading buildings list...")

	// Load buildingLocations
	buildings, err := s.db.GetBuildings()
	if err != nil {
		return fmt.Errorf("failed to init buildings: %w", err)
	}

	for _, building := range buildings {
		s.buildings[building.ID] = building
	}

	s.log.Infof("Loaded %d buildings", len(s.buildings))
	return nil
}

func (s *SimpleLogic) GetMap(_ *PlayerSession, request *rpc.GetMapRequest) (*rpc.GetMapResponse, model.Error) {
	s.log.WithField("location", request.GetLocation()).
		WithField("sessionID", request.GetSessionID()).
		Infof("GetMap request")

	return &rpc.GetMapResponse{Map: &s.GameMap}, nil
}

func (s *SimpleLogic) SelectCharacter(session *PlayerSession, request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error) {
	s.log.WithField("characterID", request.GetCharacterID()).
		WithField("sessionID", request.GetSessionID()).
		Infof("SelectCharacter request")

	char, err := s.db.GetCharacter(int(request.GetCharacterID()))
	if err != nil {
		return nil, model.ErrCharacterNotFound
	}

	session.SelectedCharacter = &char
	s.log.WithFields(logrus.Fields{
		"sessionID": request.GetSessionID(),
		"character": char,
	}).Info("User selected character")

	return &rpc.SelectCharacterResponse{}, nil
}
