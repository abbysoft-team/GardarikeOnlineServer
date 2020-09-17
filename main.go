package main

import (
	rpc "awesomeProject/rpc/generated"
	context "context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
)

const address = "localhost:27015"

func initLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type gameServer struct {
}

func (g gameServer) GetNearMap(ctx context.Context, request *rpc.GetMapRequest) (*rpc.GetMapResponse, error) {
	panic("implement me")
}

func (g gameServer) SubscribeForEvents(request *rpc.SubscribeRequest, server rpc.GameServer_SubscribeForEventsServer) error {
	panic("implement me")
}

func main() {
	listener, err := net.Listen("tcp", "localhost:27015")
	if err != nil {
		log.Fatalf("Failed to listen %s: %v", address, err)
	}

	server := grpc.NewServer()
	rpc.RegisterGameServerServer(server, &gameServer{})

	log.Printf("Server started at %s", address)
	log.Fatal(server.Serve(listener))
}
