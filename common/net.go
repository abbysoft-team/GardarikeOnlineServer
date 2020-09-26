package common

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

func WriteToSocket(msg proto.Message, address *net.UDPAddr, socket *net.UDPConn) error {
	responseData, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to serialize msg: %w", err)
	}

	if address != nil {
		_, err = socket.WriteToUDP(responseData, address)
	} else {
		_, err = socket.Write(responseData)
	}

	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}
