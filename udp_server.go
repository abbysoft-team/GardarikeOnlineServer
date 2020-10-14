package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"net"
	"projectx-server/game"
	rpc "projectx-server/rpc/generated"
)

type UDPServer struct {
	logger  *log.Entry
	socket  *net.UDPConn
	logic   game.Logic
	config  UDPServerConfig
	handler game.PacketHandler
}

type UDPServerConfig struct {
	Address        string
	ReadBufferSize int
}

func NewUDPServer(config UDPServerConfig, logic game.Logic, handler game.PacketHandler) (*UDPServer, error) {
	logger := log.WithField("module", "UDPServer")

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

	return &UDPServer{
		logger:  logger,
		socket:  conn,
		logic:   logic,
		config:  config,
		handler: handler,
	}, nil
}

func (s *UDPServer) Serve() {
	defer s.socket.Close()

	s.logger.Infof("Listen packets at %s", s.config.Address)

	var buffer [1024]byte
	for {
		bytesRead, address, err := s.socket.ReadFromUDP(buffer[0:])
		if err != nil {
			s.logger.Errorf("Read from %s failed: %v", address.String(), err)
			continue
		}

		response, err := s.handler.HandleClientPacket(buffer[0:bytesRead])
		if err != nil {
			s.logger.Errorf("Failed to handle packet: %v", err)
			continue
		}

		if response.Data != nil {
			if err, packetsSent := WriteResponse(response, address, s.socket); err != nil {
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
}

const (
	MaxPacketSize = 1024
)

// WriteResponse - writes proto message to the socket. If the specific address is not nil then
// WriteToUDP method will be used, otherwise Write method will be used.
// If the message length is more than MaxPacketSize then it will be split into packages of that size
// and sent one after another. In case of multiple parts the special MultipartResponse message
// will be sent first to inform the client about the number of packets.
// Return optional error and number of packets that have been sent (excluding MultipartResponse packet)
func WriteResponse(msg proto.Message, address *net.UDPAddr, socket *net.UDPConn) (error, int) {
	responseData, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize msg: %w", err), 0
	}

	dataLen := len(responseData)
	packetsLeft := (dataLen / MaxPacketSize) + 1

	// Write multipart response if more than one packet needed
	if packetsLeft > 1 {
		var multipartResponse rpc.Response
		multipartResponse.Data = &rpc.Response_MultipartResponse{
			MultipartResponse: &rpc.MultipartResponse{
				Parts: int64(packetsLeft),
			},
		}
		if err, _ := WriteResponse(&multipartResponse, address, socket); err != nil {
			return fmt.Errorf("failed to send multipart response message: %v", err), 0
		}
	}

	for packetsLeft > 0 {
		var dataToWrite []byte

		// Write left bytes
		if packetsLeft == 1 {
			dataToWrite = responseData[:]
		} else {
			dataToWrite = responseData[:MaxPacketSize]
		}

		if address != nil {
			_, err = socket.WriteToUDP(dataToWrite, address)
		} else {
			_, err = socket.Write(dataToWrite)
		}

		packetsLeft--
		responseData = responseData[len(dataToWrite):]
	}

	if err != nil {
		return fmt.Errorf("failed to write response: %w", err), 0
	}

	return nil, (dataLen / MaxPacketSize) + 1
}
