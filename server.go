package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"net"
	"projectx-server/common"
	"projectx-server/logic"
	rpc "projectx-server/rpc/generated"
)

type Server struct {
	logger *log.Entry
	socket *net.UDPConn
	logic  logic.Logic
	config Config
}

type Config struct {
	Address         string
	ReadBufferSize  int
	WriteBufferSize int
}

func NewServer(config Config, logic logic.Logic) (*Server, error) {
	logger := log.WithField("module", "Server")

	udpAddr, err := net.ResolveUDPAddr("udp", config.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve address: %w", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen address: %w", err)
	}

	if err := conn.SetReadBuffer(config.ReadBufferSize); err != nil {
		return nil, fmt.Errorf("failed to reserve read buffer of size %d: %w", config.ReadBufferSize, err)
	}

	if err := conn.SetWriteBuffer(config.WriteBufferSize); err != nil {
		return nil, fmt.Errorf("failed to reserve write buffer of size %d: %w", config.ReadBufferSize, err)
	}

	return &Server{
		logger: logger,
		socket: conn,
		logic:  logic,
		config: config,
	}, nil
}

func (s *Server) Serve() {
	defer s.socket.Close()

	s.logger.Infof("Listen packets at %s", s.config.Address)

	var buffer [1024]byte
	for {
		bytesRead, address, err := s.socket.ReadFromUDP(buffer[0:])
		if err != nil {
			s.logger.Errorf("Read from %s failed: %v", address.String(), err)
		}

		go s.handleClientPacket(buffer[0:bytesRead], address)
	}
}

func (s *Server) handleClientPacket(data []byte, address *net.UDPAddr) {
	var request rpc.Request
	if err := proto.Unmarshal(data, &request); err != nil {
		s.logger.WithError(err).Errorf("Failed to unmarshal client request")
		return
	}

	s.logger.WithFields(log.Fields{
		"address": address.String(),
		"request": &request,
	}).Debug("Client request")

	var requestErr error
	var response rpc.Response

	if request.GetGetMapRequest() != nil {
		getMapResponse, err := s.logic.GetMap(request.GetGetMapRequest())
		requestErr = err
		response.Data = &rpc.Response_GetMapResponse{GetMapResponse: getMapResponse}
	} else if request.GetLoginRequest() != nil {
		loginResponse, err := s.logic.Login(request.GetLoginRequest())
		requestErr = err
		response.Data = &rpc.Response_LoginResponse{LoginResponse: loginResponse}
	}

	if requestErr != nil {
		response.Data = &rpc.Response_ErrorResponse{
			ErrorResponse: &rpc.ErrorResponse{Message: requestErr.Error()},
		}
	}

	if response.Data != nil {
		if err, packetsSent := common.WriteResponse(&response, address, s.socket); err != nil {
			s.logger.
				WithError(err).
				WithField("client", address.String()).
				Error("Failed to write response to the client")
		} else {
			s.logger.
				WithField("client", address.String()).
				Infof("%d packets sent to the client", packetsSent)
		}
	}
}
