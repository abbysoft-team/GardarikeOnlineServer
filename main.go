package main

import (
	"awesomeProject/logic"
	log "github.com/sirupsen/logrus"
	"os"
)

const address = "localhost:27015"
const clientAddress = "localhost:27016"

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
	}, &logic.SimpleHandler{})
	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", address)
	}

	go startTestClient(address, clientAddress)
	server.Serve()
}

func startTestClient(serverAddress string, clientAddress string) {
	client, err := NewClient(ClientConfig{
		ListenAddress: clientAddress,
		ServerAddress: serverAddress,
	})
	if err != nil {
		log.WithError(err).Error("Failed to start client")
	}

	client.Serve()
}
