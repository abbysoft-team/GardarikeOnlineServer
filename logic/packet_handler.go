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

// TODO: refactor this method (too complex)
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
				Code:    rpc.Error(model.ErrBadRequest.GetCode()),
			},
		}

		return &response
	}

	requestName := strings.Split(fmt.Sprintf("%T", request.Data), "_")[1]

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
	} else if request.GetGetWorldMapRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetWorldMap(s, request.GetGetWorldMapRequest())
			return rpc.Response{
				Data: &rpc.Response_GetWorldMapResponse{
					GetWorldMapResponse: response,
				},
			}, err
		}

		characterRequired = false
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
	} else if request.GetGetWorkDistributionRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetWorkDistribution(s, request.GetGetWorkDistributionRequest())
			return rpc.Response{
				Data: &rpc.Response_GetWorkDistributionResponse{
					GetWorkDistributionResponse: response,
				},
			}, err
		}
	} else if request.GetCreateAccountRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.CreateAccount(request.GetCreateAccountRequest())
			return rpc.Response{
				Data: &rpc.Response_CreateAccountResponse{
					CreateAccountResponse: response,
				},
			}, err
		}

		authorizationRequired = false
		characterRequired = false
	} else if request.GetCreateCharacterRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.CreateCharacter(s, request.GetCreateCharacterRequest())
			return rpc.Response{
				Data: &rpc.Response_CreateCharacterResponse{
					CreateCharacterResponse: response,
				},
			}, err
		}
		characterRequired = false
	} else if request.GetGetResourcesRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.GetResources(s, request.GetGetResourcesRequest())
			return rpc.Response{
				Data: &rpc.Response_GetResourcesResponse{
					GetResourcesResponse: response,
				},
			}, err
		}
	} else if request.GetPlaceTownRequest() != nil {
		handleFunc = func(s *PlayerSession, r rpc.Request) (rpc.Response, model.Error) {
			response, err := p.logic.PlaceTown(s, request.GetPlaceTownRequest())
			return rpc.Response{
				Data: &rpc.Response_PlaceTownResponse{
					PlaceTownResponse: response,
				},
			}, err
		}

		characterRequired = true
		authorizationRequired = true
	} else {
		requestErr = model.ErrBadRequest
	}

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

			tx, err := p.logic.db.BeginTransaction(false, true)
			if err != nil {
				p.log.WithError(err).Error("Failed to start transaction")
				requestErr = model.ErrInternalServerError
			}

			session.Tx = tx
		}

		if requestErr == nil {
			response, requestErr = handleFunc(session, request)
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
