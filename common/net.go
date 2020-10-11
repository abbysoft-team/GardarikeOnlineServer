package common

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	rpc "projectx-server/rpc/generated"
)

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
