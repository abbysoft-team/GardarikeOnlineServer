package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"projectx-server/model"
	"projectx-server/model/postgres"
	rpc "projectx-server/rpc/generated"
)

const (
	mapChunkSize = 100
)

type Logic interface {
	GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, model.Error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, model.Error)
	SelectCharacter(request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error)
	PlaceBuilding(request *rpc.PlaceBuildingRequest) (*rpc.PlaceBuildingResponse, model.Error)
}

type SimpleLogic struct {
	gameMap   rpc.Map
	buildings map[int]model.Building
	db        model.Database
	log       *logrus.Entry
	sessions  map[string]*PlayerSession
}

func NewLogic(generator TerrainGenerator) (*SimpleLogic, error) {
	width := mapChunkSize
	height := mapChunkSize

	terrain := generator.GenerateTerrain(mapChunkSize, mapChunkSize)
	database, err := postgres.NewDatabase(postgres.Config{
		Port:      5432,
		Host:      "localhost",
		User:      "admin",
		Password:  "admin",
		DBName:    "game",
		EnableSSL: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	logic := &SimpleLogic{
		gameMap: rpc.Map{
			Width:  int32(width),
			Height: int32(height),
			Points: terrain,
		},
		buildings: make(map[int]model.Building),
		db:        database,
		log:       logrus.WithField("module", "logic"),
		sessions:  make(map[string]*PlayerSession),
	}

	logic.log.Info("Loading data from the database")
	if err := logic.load(); err != nil {
		return nil, fmt.Errorf("failed to load data from the DB: %w", err)
	}
	logic.log.Info("Logic initialization is done")

	return logic, nil
}

func (s *SimpleLogic) load() error {
	s.log.Info("Loading building locations...")

	// Load building locations
	buildingLocations, err := s.db.GetBuildingLocations()
	if err != nil {
		return fmt.Errorf("failed to load building locations: %w", err)
	}

	var rpcBuildings []*rpc.Building
	for _, building := range buildingLocations {
		rpcBuildings = append(rpcBuildings, building.ToRPC())
	}

	s.gameMap.Buildings = rpcBuildings

	s.log.Info("Loaded %d buildings on the map", len(s.gameMap.Buildings))
	s.log.Info("Loading buildings list...")

	// Load buildingLocations
	buildings, err := s.db.GetBuildings()
	if err != nil {
		return fmt.Errorf("failed to load buildings: %w", err)
	}

	for _, building := range buildings {
		s.buildings[building.ID] = building
	}

	s.log.Info("Loaded %d buildings", len(s.buildings))
	return nil
}

func (s *SimpleLogic) GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, model.Error) {
	s.log.WithField("location", request.GetLocation()).
		WithField("sessionID", request.GetSessionID()).
		Debugf("GetMap request")

	_, authorized := s.sessions[request.GetSessionID()]
	if !authorized {
		return nil, model.ErrNotAuthorized
	}

	return &rpc.GetMapResponse{Map: &s.gameMap}, nil
}

func (s *SimpleLogic) SelectCharacter(request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, model.Error) {
	s.log.WithField("characterID", request.GetCharacterID()).
		WithField("sessionID", request.GetSessionID()).
		Debugf("SelectCharacter request")

	session, authorized := s.sessions[request.GetSessionID()]
	if !authorized {
		return nil, model.ErrNotAuthorized
	}

	char, err := s.db.GetCharacter(int(request.GetCharacterID()))
	if err != nil {
		return nil, model.ErrCharacterNotFound
	}

	session.SelectedCharacter = &char
	s.log.WithFields(logrus.Fields{
		"accountID": request.GetSessionID(),
		"character": char,
	}).Info("User selected character")

	return &rpc.SelectCharacterResponse{}, nil
}
