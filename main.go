package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"projectx-server/game"
)

const (
	version = "0.0.4"
)

func initLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please, specify listen address in first argument")
	}

	if os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	initLogging()
	log.Printf("ProjectX UDPServer v%s", version)

	config := Config{
		Endpoint: os.Args[1],
	}

	gameLogic := game.NewLogic(game.NewSimplexMapGenerator(5, 1.5))
	handler := game.NewPacketHandler(gameLogic)

	//server, err := NewUDPServer(config, game, )
	server, err := NewServer(config, gameLogic, handler)
	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", os.Args[1])
	}

	log.Fatal(server.Serve())
}
