package server

import (
	logic2 "abbysoft/gardarike-online/logic"
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

type PacketHandler struct {
	logic logic2.Logic
	log   *logrus.Entry
}

func NewPacketHandler(logic logic2.Logic) PacketHandler {
	return PacketHandler{
		logic: logic,
		log:   logrus.WithField("module", "packet_handler"),
	}
}

func (p *PacketHandler) HandleClientPacket(data []byte) *rpc.Response {
	var request rpc.Request
	var requestErr model.Error
	var response rpc.Response

	if err := proto.Unmarshal(data, &request); err != nil || len(data) == 0 {
		p.log.WithError(err).Error("Failed to serialize client request")

		response.Data = &rpc.Response_ErrorResponse{
			ErrorResponse: &rpc.ErrorResponse{
				Message: model.ErrBadRequest.GetMessage(),
				Code:    int64(model.ErrBadRequest.GetCode()),
			},
		}

		return &response
	}

	if request.GetLoginRequest() != nil {
		loginResponse, err := p.logic.Login(request.GetLoginRequest())
		requestErr = err
		response.Data = &rpc.Response_LoginResponse{LoginResponse: loginResponse}
	} else if request.GetGetMapRequest() != nil {
		getMapResponse, err := p.logic.GetMap(request.GetGetMapRequest())
		requestErr = err
		response.Data = &rpc.Response_GetMapResponse{GetMapResponse: getMapResponse}
	} else if request.GetSelectCharacterRequest() != nil {
		selectCharResponse, err := p.logic.SelectCharacter(request.GetSelectCharacterRequest())
		requestErr = err
		response.Data = &rpc.Response_SelectCharacterResponse{SelectCharacterResponse: selectCharResponse}
	} else if request.GetPlaceBuildingRequest() != nil {
		placeBuildingResponse, err := p.logic.PlaceBuilding(request.GetPlaceBuildingRequest())
		requestErr = err
		response.Data = &rpc.Response_PlaceBuildingResponse{PlaceBuildingResponse: placeBuildingResponse}
	} else if request.GetSendChatMessageRequest() != nil {
		sendChatMessageResponse, err := p.logic.SendChatMessage(request.GetSendChatMessageRequest())
		requestErr = err
		response.Data = &rpc.Response_SendChatMessageResponse{SendChatMessageResponse: sendChatMessageResponse}
	}

	if requestErr != nil {
		response.Data = &rpc.Response_ErrorResponse{
			ErrorResponse: &rpc.ErrorResponse{
				Message: requestErr.GetMessage(),
				Code:    int64(requestErr.GetCode()),
			},
		}
	}

	return &response
}
