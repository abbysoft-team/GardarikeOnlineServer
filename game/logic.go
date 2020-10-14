package game

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

type SimpleLogic struct {
	gameMap rpc.Map
}

func NewLogic(generator TerrainGenerator) *SimpleLogic {
	width := mapChunkSize
	height := mapChunkSize

	terrain := generator.GenerateTerrain(mapChunkSize, mapChunkSize)
	return &SimpleLogic{gameMap: rpc.Map{
		Width:  int32(width),
		Height: int32(height),
		Points: terrain,
	}}
}

func (s *SimpleLogic) GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error) {
	return &rpc.GetMapResponse{Map: &s.gameMap}, nil
}

func (s *SimpleLogic) Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	return &rpc.LoginResponse{}, nil
}
