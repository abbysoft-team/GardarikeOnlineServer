package tests

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

const (
	//serverEndpoint = "tcp://89.108.99.2:27015"
	serverEndpoint = "tcp://localhost:8500"
	//serverEventEndpoint = "tcp://89.108.99.2:27016"
	serverEventEndpoint = "tcp://localhost:8501"
	requestTimeout      = 1 * time.Second
)

var client *Client
var sessionID string

func TestMain(m *testing.M) {
	testClient, err := NewClient(ClientConfig{
		ServerEndpoint:      serverEndpoint,
		ServerEventEndpoint: serverEventEndpoint,
		RequestTimeout:      requestTimeout,
	})
	if err != nil {
		log.Fatalf("Failed to init test client: %v", err)
		os.Exit(1)
	}

	client = testClient
	os.Exit(m.Run())
}
