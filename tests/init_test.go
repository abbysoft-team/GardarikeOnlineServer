// +build !remote_tests

package tests

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

const (
	//serverEndpoint = "tcp://89.108.99.2:27015"
	serverEndpoint = "tcp://localhost:27015"
	//serverEventEndpoint = "tcp://89.108.99.2:27016"
	serverEventEndpoint = "tcp://localhost:27016"
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

	if len(os.Getenv("RUN_REMOTE_TESTS")) > 0 {
		os.Exit(m.Run())
	} else {
		log.Info("Skipping remote tests because RUN_REMOTE_TESTS isn't specified")
	}
}
