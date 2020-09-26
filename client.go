package main

import (
	"awesomeProject/common"
	rpc "awesomeProject/rpc/generated"
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Client struct {
	socket        *net.UDPConn
	listenAddress *net.UDPAddr
	serverAddress *net.UDPAddr
	logger        *log.Entry
}

type ClientConfig struct {
	ListenAddress string
	ServerAddress string
}

func NewClient(config ClientConfig) (*Client, error) {
	listenAddress, err := net.ResolveUDPAddr("udp", config.ListenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve listen address: %w", err)
	}

	serverAddress, err := net.ResolveUDPAddr("udp", config.ServerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve server address: %w", err)
	}

	conn, err := net.DialUDP("udp", listenAddress, serverAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create client socket: %w", err)
	}

	return &Client{
		socket:        conn,
		listenAddress: listenAddress,
		serverAddress: serverAddress,
		logger:        log.WithField("module", "Client"),
	}, nil
}

func (c *Client) Serve() {
	defer c.socket.Close()

	go c.readResponses()
	for {
		time.Sleep(1 * time.Second)

		var request rpc.Request
		request.Data = &rpc.Request_GetMapRequest{
			GetMapRequest: &rpc.GetMapRequest{
				Location: &rpc.Vector3D{
					X: 10,
					Y: 15,
					Z: 20,
				},
			},
		}

		if err := common.WriteToSocket(&request, nil, c.socket); err != nil {
			c.logger.WithError(err).Error("Failed to write request")
		}
	}
}

func (c *Client) readResponses() {
	var buffer [1024]byte
	for {
		bytesRead, err := c.socket.Read(buffer[0:])
		if err != nil {
			c.logger.WithError(err).Error("Failed to read response from the server")
			continue
		}

		var response rpc.Response
		if err := proto.Unmarshal(buffer[0:bytesRead], &response); err != nil {
			c.logger.WithError(err).Error("Failed to deserialize server response")
		}

		c.logger.
			WithField("response", response.Data).
			Infof("Server respond with %d bytes", bytesRead)
	}
}
