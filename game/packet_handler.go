package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	rpc "projectx-server/rpc/generated"
)

type PacketHandler struct {
	logic Logic
}

func NewPacketHandler(logic Logic) PacketHandler {
	return PacketHandler{
		logic: logic,
	}
}

func (p *PacketHandler) HandleClientPacket(data []byte) (*rpc.Response, error) {
	var request rpc.Request
	if err := proto.Unmarshal(data, &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal packet: %w", err)
	}

	var requestErr error
	var response rpc.Response

	if request.GetGetMapRequest() != nil {
		getMapResponse, err := p.logic.GetMap(request.GetGetMapRequest())
		requestErr = err
		response.Data = &rpc.Response_GetMapResponse{GetMapResponse: getMapResponse}
	} else if request.GetLoginRequest() != nil {
		loginResponse, err := p.logic.Login(request.GetLoginRequest())
		requestErr = err
		response.Data = &rpc.Response_LoginResponse{LoginResponse: loginResponse}
	} else if request.GetSelectCharacterRequest() != nil {
		selectCharResponse, err := p.logic.SelectCharacter(request.GetSelectCharacterRequest())
		requestErr = err
		response.Data = &rpc.Response_GetCharacterResponse{GetCharacterResponse: selectCharResponse}
	}

	if requestErr != nil {
		response.Data = &rpc.Response_ErrorResponse{
			ErrorResponse: &rpc.ErrorResponse{Message: requestErr.Error()},
		}
	}

	return &response, nil
}
