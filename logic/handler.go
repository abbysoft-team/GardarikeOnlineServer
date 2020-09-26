package logic

import (
	rpc "awesomeProject/rpc/generated"
	"fmt"
)

type Handler interface {
	GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error)
	Subscribe(request *rpc.SubscribeRequest) error
}

type SimpleHandler struct {
}

func (s SimpleHandler) GetMap(request *rpc.GetMapRequest) (*rpc.GetMapResponse, error) {
	return nil, fmt.Errorf("getMap is unimplemented")
}

func (s SimpleHandler) Subscribe(request *rpc.SubscribeRequest) error {
	return fmt.Errorf("subscribe is unimplemented")
}
