package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	zmq "github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
	"projectx-server/game"
)

type Server struct {
	context *zmq.Context
	socket  *zmq.Socket
	config  Config
	log     *log.Entry
	logic   game.Logic
	handler game.PacketHandler
}

type Config struct {
	Endpoint string // ZMQ endpoint string (e.g. tcp://*:555)
}

func NewServer(config Config, logic game.Logic, handler game.PacketHandler) (*Server, error) {
	context, err := zmq.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create zmq context: %w", err)
	}

	sock, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		return nil, fmt.Errorf("failed to create ZMQ REP socket: %w", err)
	}

	logger := log.WithField("module", "server")

	return &Server{
		socket:  sock,
		config:  config,
		log:     logger,
		logic:   logic,
		handler: handler,
		context: context,
	}, nil
}

func (s *Server) Serve() error {
	if err := s.socket.Bind(s.config.Endpoint); err != nil {
		return fmt.Errorf("failed to bind server socket to address %s: %w", s.config.Endpoint, err)
	}

	s.log.Infof("Logic listen on %s", s.config.Endpoint)

	for {
		packet, err := s.socket.Recv(0)
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

		if _, err := s.socket.Send(string(respBytes), 0); err != nil {
			s.log.Errorf("Failed to send answer to the client: %v", err)
		}
	}
}
