package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/sirupsen/logrus"
)

func (s *SimpleLogic) GetChatHistory(session *PlayerSession, request *rpc.GetChatHistoryRequest) (*rpc.GetChatHistoryResponse, model.Error) {
	s.log.WithFields(logrus.Fields{
		"sessionID": request.SessionID,
		"offset":    request.Offset,
		"count":     request.Count,
	}).Info("GetChatHistory")

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

	var rpcMessages []*rpc.ChatMessage
	for _, message := range messages {
		rpcMessages = append(rpcMessages, message.ToRPC())
	}

	return &rpc.GetChatHistoryResponse{
		Messages: rpcMessages,
	}, nil
}
