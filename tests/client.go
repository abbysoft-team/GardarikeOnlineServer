package tests

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	socket      *zmq.Socket
	eventSocket *zmq.Socket
	context     *zmq.Context
	logger      *log.Entry
	config      ClientConfig
	eventChan   chan *rpc.Event
}

type ClientConfig struct {
	ServerEndpoint      string
	ServerEventEndpoint string
	RequestTimeout      time.Duration
}

func NewClient(config ClientConfig) (*Client, error) {
	context, err := zmq.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create zmq context: %w", err)
	}

	socket, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		return nil, fmt.Errorf("failed to create zmq socket: %w", err)
	}

	if err = socket.SetSndtimeo(config.RequestTimeout); err != nil {
		return nil, fmt.Errorf("failed to set socket send timeout option: %w", err)
	}

	if err = socket.SetRcvtimeo(config.RequestTimeout); err != nil {
		return nil, fmt.Errorf("failed to set socket recv timeout option: %w", err)
	}

	if err = socket.Connect(config.ServerEndpoint); err != nil {
		return nil, fmt.Errorf("failed to connect to the server: %w", err)
	}

	eventSocket, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe for server events: %w", err)
	}

	logger := log.WithField("module", "Client")

	if err = eventSocket.SetSubscribe("GLOB"); err != nil {
		return nil, fmt.Errorf("failed to subscribe to GLOB channel: %w", err)
	}

	if err = eventSocket.Connect(config.ServerEventEndpoint); err != nil {
		return nil, fmt.Errorf("failed to subscribe to the server events socket: %w", err)
	}

	client := &Client{
		socket:      socket,
		eventSocket: eventSocket,
		logger:      logger,
		config:      config,
		context:     context,
		eventChan:   make(chan *rpc.Event, 10),
	}

	go client.pollEvents()

	return client, nil
}

func (c *Client) pollEvents() {
	for {
		event, err := c.pollEvent()
		if err != nil {
			c.logger.Error("Failed to poll event: %w", err)
			time.Sleep(1 * time.Second)
			continue
		}

		c.eventChan <- event
	}
}

func (c *Client) pollEvent() (*rpc.Event, error) {
	eventParts, err := c.eventSocket.RecvMessage(0)
	if err != nil {
		return nil, fmt.Errorf("failed to recv event: %w", err)
	}

	var event rpc.Event
	if err = proto.Unmarshal([]byte(eventParts[1]), &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}

	return &event, nil
}

func (c *Client) SendMessage(message proto.Message) {
	requestBytes, err := proto.Marshal(message)
	if err != nil {
		c.logger.WithError(err).Error("Failed to marshal request")
		return
	}

	if _, err := c.socket.Send(string(requestBytes), zmq.DONTWAIT); err != nil {
		c.logger.WithError(err).Error("Failed to send request to the server")
	}
}

func (c *Client) SendRequest(request rpc.Request) (*rpc.Response, error) {
	requestBytes, err := proto.Marshal(&request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	c.logger.WithFields(log.Fields{
		"server":  c.config.ServerEndpoint,
		"request": fmt.Sprintf("%T", request.Data),
		"bytes":   fmt.Sprintf("%x", requestBytes),
	}).Info("Send request to the server")

	if _, err := c.socket.Send(string(requestBytes), zmq.DONTWAIT); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	c.logger.Printf("%T sent to the server", request.Data)

	if response, err := c.readResponse(); err != nil {
		return nil, fmt.Errorf("failed to read response to the server: %w", err)
	} else {
		if errorResp := response.GetErrorResponse(); errorResp != nil {
			return nil, model.NewError(errorResp.Message, errorResp.Code)
		}
		return response, nil
	}
}

func (c *Client) readResponse() (*rpc.Response, error) {
	bytesRead, err := c.socket.Recv(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from the server: %w", err)
	}

	var response rpc.Response
	if err := proto.Unmarshal([]byte(bytesRead), &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal server response: %w", err)
	}

	if response.GetMultipartResponse() != nil {
		return nil, fmt.Errorf("multipart response received")
	}

	c.logger.
		WithField("response", response.Data).
		Infof("Server respond with %d bytes", len(bytesRead))
	return &response, nil
}
