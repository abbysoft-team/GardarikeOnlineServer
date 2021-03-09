package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
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

type handleFunc func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error)

type requestHandler struct {
	handleFunc            handleFunc
	authorizationRequired bool
	characterRequired     bool
}

func (p *PacketHandler) getHandleFunc(request rpc.Request) *requestHandler {
	var handler requestHandler

	if request.GetLoginRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.Login(request.GetLoginRequest())
			return rpc.Response{
				Data: &rpc.Response_LoginResponse{
					LoginResponse: response,
				},
			}, err
		}
		handler.authorizationRequired = false
		handler.characterRequired = false
	} else if request.GetGetWorldMapRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetWorldMap(s, request.GetGetWorldMapRequest())
			return rpc.Response{
				Data: &rpc.Response_GetWorldMapResponse{
					GetWorldMapResponse: response,
				},
			}, err
		}

		handler.characterRequired = false
	} else if request.GetSelectCharacterRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.SelectCharacter(s, request.GetSelectCharacterRequest())
			return rpc.Response{
				Data: &rpc.Response_SelectCharacterResponse{
					SelectCharacterResponse: response,
				},
			}, err
		}
		handler.characterRequired = false
	} else if request.GetSendChatMessageRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.SendChatMessage(s, request.GetSendChatMessageRequest())
			return rpc.Response{
				Data: &rpc.Response_SendChatMessageResponse{
					SendChatMessageResponse: response,
				},
			}, err
		}
	} else if request.GetGetChatHistoryRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetChatHistory(s, request.GetGetChatHistoryRequest())
			return rpc.Response{
				Data: &rpc.Response_GetChatHistoryResponse{
					GetChatHistoryResponse: response,
				},
			}, err
		}
	} else if request.GetGetWorkDistributionRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetWorkDistribution(s, request.GetGetWorkDistributionRequest())
			return rpc.Response{
				Data: &rpc.Response_GetWorkDistributionResponse{
					GetWorkDistributionResponse: response,
				},
			}, err
		}
	} else if request.GetCreateAccountRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.CreateAccount(request.GetCreateAccountRequest())
			return rpc.Response{
				Data: &rpc.Response_CreateAccountResponse{
					CreateAccountResponse: response,
				},
			}, err
		}

		handler.authorizationRequired = false
		handler.characterRequired = false
	} else if request.GetCreateCharacterRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.CreateCharacter(s, request.GetCreateCharacterRequest())
			return rpc.Response{
				Data: &rpc.Response_CreateCharacterResponse{
					CreateCharacterResponse: response,
				},
			}, err
		}
		handler.characterRequired = false
	} else if request.GetGetResourcesRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetResources(s, request.GetGetResourcesRequest())
			return rpc.Response{
				Data: &rpc.Response_GetResourcesResponse{
					GetResourcesResponse: response,
				},
			}, err
		}
	} else if request.GetPlaceTownRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.PlaceTown(s, request.GetPlaceTownRequest())
			return rpc.Response{
				Data: &rpc.Response_PlaceTownResponse{
					PlaceTownResponse: response,
				},
			}, err
		}

		handler.characterRequired = true
		handler.authorizationRequired = true
	} else if request.GetPlaceBuildingRequest() != nil {
		handler.handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.PlaceBuilding(s, r.GetPlaceBuildingRequest())
			return rpc.Response{
				Data: &rpc.Response_PlaceBuildingResponse{
					PlaceBuildingResponse: response,
				},
			}, err
		}
	} else {
		return nil
	}

	return &handler
}

// TODO: refactor this method (too complex)
func (p *PacketHandler) HandleClientPacket(data []byte) *rpc.Response {
	var request rpc.Request
	var requestErr model.Error
	var response rpc.Response

	if err := proto.Unmarshal(data, &request); err != nil || len(data) == 0 {
		p.log.WithError(err).Error("Failed to serialize client request")

		response.Data = &rpc.Response_ErrorResponse{
			ErrorResponse: &rpc.ErrorResponse{
				Message: model.ErrBadRequest.GetMessage(),
				Code:    rpc.Error(model.ErrBadRequest.GetCode()),
			},
		}

		return &response
	}

	requestName := strings.Split(fmt.Sprintf("%T", request.Data), "_")[1]

	var sessionID string
	var authorized bool
	var session *PlayerSession

	// Check session
	sessionRegexp, _ := regexp.Compile("sessionID:\"([^\"]*)\"")
	sessionSubmatch := sessionRegexp.FindStringSubmatch(request.String())
	if len(sessionSubmatch) == 2 {
		sessionID = sessionSubmatch[1]
		session, authorized = p.logic.sessions[sessionID]

		if session != nil {
			session.LastRequestTime = time.Now()
		}
	}

	handler := p.getHandleFunc(request)
	if handler == nil {
		requestErr = model.ErrBadRequest
	}

	if handler != nil && !authorized && handler.authorizationRequired {
		requestErr = model.ErrNotAuthorized
	} else if handler != nil &&
		session != nil &&
		requestErr == nil &&
		handler.characterRequired &&
		session.SelectedCharacter == nil {
		requestErr = model.ErrCharacterNotSelected
	}

	if handler != nil && handler.handleFunc != nil && requestErr == nil {
		if session != nil {
			session.Mutex.Lock()

			tx, err := p.logic.db.BeginTransaction(false, true)
			if err != nil {
				p.log.WithError(err).Error("Failed to start transaction")
				requestErr = model.ErrInternalServerError
			}

			session.Tx = tx
		}

		if requestErr == nil {
			response, requestErr = handler.handleFunc(session, request)
		}
		if session != nil {
			// Only commit should be handled, rollback is happened automatically on errors
			if session.Tx != nil && !session.Tx.IsCompleted() {
				if err := session.Tx.EndTransaction(); err != nil {
					p.log.WithError(err).Error("Failed to commit transaction")
					requestErr = model.ErrInternalServerError
				}
			}

			session.Mutex.Unlock()
		}
	}

	if requestErr != nil {
		p.log.WithField("requestName", requestName).Infof("Sending error response: %v", requestErr.Error())

		response = rpc.Response{
			Data: &rpc.Response_ErrorResponse{
				ErrorResponse: &rpc.ErrorResponse{
					Message: requestErr.GetMessage(),
					Code:    rpc.Error(requestErr.GetCode()),
				},
			},
		}
	}

	return &response
}
