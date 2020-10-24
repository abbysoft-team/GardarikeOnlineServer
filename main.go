package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	version = "0.1.2"
)

func initLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func printUsage() {
	fmt.Printf("Usage: gardarike requestEndpoint eventEndpoint")
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println(version)
		os.Exit(0)
	}

	if len(os.Args) < 3 {
		printUsage()
		os.Exit(1)
	}

	initLogging()
	log.Printf("ProjectX UDPServer v%s", version)

	config := Config{
		RequestEndpoint: os.Args[1],
		EventEndpoint:   os.Args[2],
	}

	//server, err := NewUDPServer(config, game, )
	server, err := NewServer(config)
	if err != nil {
		log.WithError(err).Fatalf("Failed to start server on %s", os.Args[1])
	}

	log.Fatal(server.Serve())
}
