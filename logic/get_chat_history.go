package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/sirupsen/logrus"
)

func toRpcMessages(messages []model.ChatMessage) (result []*rpc.ChatMessage) {
	for _, message := range messages {
		result = append(result, message.ToRPC())
	}

	return
}

func (s *SimpleLogic) GetChatHistory(request *rpc.GetChatHistoryRequest) (*rpc.GetChatHistoryResponse, model.Error) {
	s.log.WithFields(logrus.Fields{
		"sessionID": request.SessionID,
		"offset":    request.Offset,
		"count":     request.Count,
	}).Info("GetChatHistory")

	_, err := s.checkAuthorization(request.SessionID)
	if err != nil {
		return nil, model.ErrNotAuthorized
	}

	offset := 0
	limit := 10
	if request.Offset != 0 {
		offset = int(request.Offset)
	}
	if request.Count != 0 {
		limit = int(request.Count)
	}

	messages, dbErr := s.db.GetChatMessages(offset, limit)
	if dbErr != nil {
		s.log.WithError(dbErr).Error("Failed to GetChatMessages")
		return nil, model.ErrInternalServerError
	}

	return &rpc.GetChatHistoryResponse{
		Messages: toRpcMessages(messages),
	}, nil
}
