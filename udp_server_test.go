package main

import (
	"github.com/golang/protobuf/proto"
	"net"
	rpc "projectx-server/rpc/generated"
	"testing"
)

func generateStringOfLen(len int) string {
	result := make([]byte, len)
	for i := 0; i < len; i++ {
		result[i] = 'A'
	}

	result[len-1] = 'B'

	return string(result)
}

func TestWriteToSocket(t *testing.T) {
	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:11111")
	if err != nil {
		t.Fatalf("Failed to resolve test address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatalf("Failed to create test socket: %v", err)
	}

	tests := []struct {
		name          string
		requestLength float64
	}{
		{"One-packet request sent successfully", MaxPacketSize / 2},
		{"Two-packet request sent successfully", MaxPacketSize * 1.5},
		{"Three-packet request sent successfully", MaxPacketSize * 2.5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var request rpc.Request
			request.Data = &rpc.Request_LoginRequest{
				LoginRequest: &rpc.LoginRequest{
					Username: generateStringOfLen(int(test.requestLength)),
					Password: "",
				},
			}

			var buffer [MaxPacketSize]byte
			rightNumberOfPackets := int(test.requestLength/MaxPacketSize) + 1

			err, packetsSent := WriteResponse(&request, udpAddr, conn)
			if err != nil {
				t.Fatalf("WriteResponse failed: %v", err)
			} else if packetsSent != rightNumberOfPackets {
				t.Errorf("WriteResponse expect %d packet to be sent but %d packets have been sent",
					rightNumberOfPackets,
					packetsSent)
			}

			var responseBuffer []byte
			if packetsSent > 1 {
				_, err := conn.Read(buffer[0:])
				if err != nil {
					t.Errorf("Failed to read multipart packet: %v", err)
				}
			}

			for i := 0; i < packetsSent; i++ {
				if bytes, err := conn.Read(buffer[0:]); err != nil {
					t.Fatalf("Failed to read packet data from the socket: %v", err)
				} else {
					responseBuffer = append(responseBuffer, buffer[0:bytes]...)
				}
			}

			var response rpc.Request
			err = proto.Unmarshal(responseBuffer, &response)
			if err != nil {
				t.Fatalf("Failed to serialize response: %v", err)
			}
		})
	}
}
