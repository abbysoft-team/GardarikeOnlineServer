package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"projectx-server/logic"
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

	if len(os.Args) < 2 {
		log.Fatal("Please, specify listen address in first argument")
	}

	server, err := NewServer(Config{
		Address:        os.Args[1],
		ReadBufferSize: readBufferSize,
	}, logic.NewServerLogic(logic.NewSimplexMapGenerator(5, 1.5)))

	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", address)
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
