package logic

import (
	log "github.com/sirupsen/logrus"
	rpc "projectx-server/rpc/generated"
)

const (
	mapChunkSize = 100
)

type Logic interface {
	GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error)
	Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error)
}

type ServerLogic struct {
	gameMap rpc.Map
}

func NewServerLogic(generator TerrainGenerator) *ServerLogic {
	width := mapChunkSize
	height := mapChunkSize

	terrain := generator.GenerateTerrain(mapChunkSize, mapChunkSize)

	log.WithField("size", len(terrain)*4).Infof("Terrain generated")

	return &ServerLogic{gameMap: rpc.Map{
		Width:  int32(width),
		Height: int32(height),
		Points: terrain,
	}}
}

func (s *ServerLogic) GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error) {
	return &rpc.GetMapResponse{Map: &s.gameMap}, nil
}

func (s *ServerLogic) Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	return &rpc.LoginResponse{}, nil
}
