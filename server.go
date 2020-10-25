package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
	"projectx-server/game"
	"projectx-server/model/postgres"
	rpc "projectx-server/rpc/generated"
)

type Server struct {
	context     *zmq.Context
	requestSock *zmq.Socket
	eventSock   *zmq.Socket
	config      Config
	log         *log.Entry
	logic       game.Logic
	handler     game.PacketHandler
	eventsChan  chan *rpc.Event
}

type Config struct {
	RequestEndpoint string // Listens for requests on this endpoint (e.g. tcp://*:555)
	EventEndpoint   string // Publish events on this endpoint
}

func NewServer(config Config, dbConfig postgres.Config, generatorConfig game.TerrainGeneratorConfig) (*Server, error) {
	context, err := zmq.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create zmq context: %w", err)
	}

	sock, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZMQ REP request socket: %w", err)
	}

	eventSock, err := zmq.NewSocket(zmq.PUB)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZMQ PUB event socket: %w", err)
	}

	logger := log.WithField("module", "server")

	eventsChan := make(chan *rpc.Event, 10)
	gameLogic, err := game.NewLogic(
		game.NewSimplexTerrainGenerator(generatorConfig),
		eventsChan,
		dbConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to init game logic: %w", err)
	}

	handler := game.NewPacketHandler(gameLogic)

	return &Server{
		requestSock: sock,
		eventSock:   eventSock,
		config:      config,
		log:         logger,
		logic:       gameLogic,
		handler:     handler,
		context:     context,
		eventsChan:  eventsChan,
	}, nil
}

func (s *Server) publishEvent(event *rpc.Event) {
	logger := s.log.WithField("event", fmt.Sprintf("%T", event.Payload))

	bytes, err := proto.Marshal(event)
	if err != nil {
		logger.Errorf("Failed to marshal server event: %v", err)
		return
	}

	if _, err := s.eventSock.Send(string(bytes), zmq.DONTWAIT); err != nil {
		logger.WithError(err).Error("Failed to push server event: %v", err)
	} else {
		logger.WithField("payload", event.Payload).Info("Event published to the clients")
	}
}

func (s *Server) serveEvents() {
	for event := range s.eventsChan {
		s.publishEvent(event)
	}
}

func (s *Server) Serve() error {
	if err := s.requestSock.Bind(s.config.RequestEndpoint); err != nil {
		return fmt.Errorf("failed to bind server requestSock to address %s: %w", s.config.RequestEndpoint, err)
	}

	if err := s.eventSock.Bind(s.config.EventEndpoint); err != nil {
		return fmt.Errorf("failed to bind server eventSock to address: %s: %w", s.config.EventEndpoint, err)
	}

	s.log.WithFields(log.Fields{
		"requestEndpoint": s.config.RequestEndpoint,
		"eventEndpoint":   s.config.EventEndpoint,
	}).Infof("Server started")

	go s.serveEvents()

	for {
		packet, err := s.requestSock.Recv(0)
		if err != nil {
			s.log.Errorf("Failed to read client packet: %v", err)
			continue
		}

		s.log.Debugf("Read %d bytes from client", len(packet))

		resp := s.handler.HandleClientPacket([]byte(packet))

		respBytes, err := proto.Marshal(resp)
		if err != nil {
			s.log.Errorf("Failed to marshal server response: %v", err)
			continue
		}

		s.log.Infof("Sending %T response to the client (%d bytes)", resp.Data, len(respBytes))

		if _, err := s.requestSock.Send(string(respBytes), 0); err != nil {
			s.log.Errorf("Failed to send answer to the client: %v", err)
		}
	}
}
