package main

import (
	"awesomeProject/logic"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

const address = "localhost:27015"
const clientPortBase = 2015
const clientsCount = 5
const readBufferSize = 1024 * 10 // 10 KB

func initLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type gameServer struct {
}

func main() {
	log.SetLevel(log.DebugLevel)

	server, err := NewServer(Config{
		Address:        address,
		ReadBufferSize: readBufferSize,
	}, logic.NewServerLogic(logic.NewSimplexMapGenerator(5, 1.5)))

	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", address)
	}

	for i := 0; i < clientsCount; i++ {
		go startTestClient(address, fmt.Sprintf("localhost:%d", clientPortBase+i))
	}

	server.Serve()
}

func startTestClient(serverAddress string, clientAddress string) {
	client, err := NewClient(ClientConfig{
		ListenAddress: clientAddress,
		ServerAddress: serverAddress,
	})
	if err != nil {
		log.WithError(err).Error("Failed to start client")
		return
	}

	client.Serve()
}
