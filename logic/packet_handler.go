package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
)

type PacketHandler struct {
	logic *SimpleLogic
	log   *logrus.Entry
}

func NewPacketHandler(logic *SimpleLogic) PacketHandler {
	return PacketHandler{
		logic: logic,
		log:   logrus.WithField("module", "packet_handler"),
	}
}

func (p *PacketHandler) HandleClientPacket(data []byte) *rpc.Response {
	var request rpc.Request
	var requestErr model.Error
	var response rpc.Response
	authorizationRequired := true
	characterRequired := true

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

	var handleFunc func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error)

	if request.GetLoginRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.Login(request.GetLoginRequest())
			return rpc.Response{
				Data: &rpc.Response_LoginResponse{
					LoginResponse: response,
				},
			}, err
		}
		authorizationRequired = false
		characterRequired = false
	} else if request.GetGetMapRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetMap(s, request.GetGetMapRequest())
			return rpc.Response{
				Data: &rpc.Response_GetMapResponse{
					GetMapResponse: response,
				},
			}, err
		}
	} else if request.GetSelectCharacterRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.SelectCharacter(s, request.GetSelectCharacterRequest())
			return rpc.Response{
				Data: &rpc.Response_SelectCharacterResponse{
					SelectCharacterResponse: response,
				},
			}, err
		}
		characterRequired = false
	} else if request.GetPlaceBuildingRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.PlaceBuilding(s, request.GetPlaceBuildingRequest())
			return rpc.Response{
				Data: &rpc.Response_PlaceBuildingResponse{
					PlaceBuildingResponse: response,
				},
			}, err
		}
	} else if request.GetSendChatMessageRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.SendChatMessage(s, request.GetSendChatMessageRequest())
			return rpc.Response{
				Data: &rpc.Response_SendChatMessageResponse{
					SendChatMessageResponse: response,
				},
			}, err
		}
	} else if request.GetGetChatHistoryRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetChatHistory(s, request.GetGetChatHistoryRequest())
			return rpc.Response{
				Data: &rpc.Response_GetChatHistoryResponse{
					GetChatHistoryResponse: response,
				},
			}, err
		}
	} else {
		requestErr = model.ErrBadRequest
	}

	var sessionID string
	var authorized bool
	var session *PlayerSession

	// Check session
	sessionRegexp, _ := regexp.Compile("sessionID:\"(.*)\"")
	sessionSubmatch := sessionRegexp.FindStringSubmatch(request.String())
	if len(sessionSubmatch) == 2 {
		sessionID = sessionSubmatch[1]
		session, authorized = p.logic.sessions[sessionID]

		if session != nil {
			session.LastRequestTime = time.Now()
		}
	}

	if !authorized && authorizationRequired {
		requestErr = model.ErrNotAuthorized
	} else if session != nil &&
		requestErr == nil &&
		characterRequired &&
		session.SelectedCharacter == nil {
		requestErr = model.ErrCharacterNotSelected
	}

	if handleFunc != nil && requestErr == nil {
		if session != nil {
			session.Mutex.Lock()
		}

		response, requestErr = handleFunc(session, request)

		if session != nil {
			session.Mutex.Unlock()
		}
	}

	if requestErr != nil {
		response = rpc.Response{
			Data: &rpc.Response_ErrorResponse{
				ErrorResponse: &rpc.ErrorResponse{
					Message: requestErr.GetMessage(),
					Code:    int64(requestErr.GetCode()),
				},
			},
		}
	}

	return &response
}
