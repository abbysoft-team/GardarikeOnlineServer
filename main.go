package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"projectx-server/common"
	"projectx-server/logic"
)

const (
	clientPortBase = 2015
	clientsCount   = 5
	version        = "0.0.4"
)

func initLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type gameServer struct {
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please, specify listen address in first argument")
	}

	if os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	log.SetLevel(log.DebugLevel)
	log.Printf("ProjectX Server v%s", version)

	server, err := NewServer(Config{
		Address:         os.Args[1],
		ReadBufferSize:  common.MaxPacketSize,
		WriteBufferSize: common.MaxPacketSize,
	}, logic.NewServerLogic(logic.NewSimplexMapGenerator(5, 1.5)))

	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", os.Args[1])
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
