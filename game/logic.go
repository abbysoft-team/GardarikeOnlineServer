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
	GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error)
	SelectCharacter(request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, error)
}

type SimpleLogic struct {
	gameMap rpc.Map
	db      model.Database
	log     *logrus.Entry
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

	return &SimpleLogic{
		gameMap: rpc.Map{
			Width:  int32(width),
			Height: int32(height),
			Points: terrain,
		},
		db:  database,
		log: logrus.WithField("module", "logic"),
	}, nil
}

func (s *SimpleLogic) GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error) {
	s.log.WithField("location", request.GetLocation()).Debugf("GetMap request")
	return &rpc.GetMapResponse{Map: &s.gameMap}, nil
}

func (s *SimpleLogic) SelectCharacter(request *rpc.SelectCharacterRequest) (*rpc.SelectCharacterResponse, error) {
	return nil, fmt.Errorf("unimplemented")
}
