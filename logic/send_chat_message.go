package logic

import (
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) SendChatMessage(session *PlayerSession, request *rpc.SendChatMessageRequest) (*rpc.SendChatMessageResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": request.SessionID,
		"text":      request.Text,
	}).Info("SendChatMessage")

	if len(request.Text) > s.config.ChatMessageMaxLength {
		return nil, model.ErrMessageTooLong
	}

	message := model.ChatMessage{
		ID:     0,
		Sender: session.SelectedCharacter.Name,
		Text:   request.Text,
	}

	if insertedID, err := s.db.AddChatMessage(message); err != nil {
		s.log.WithError(err).Error("Failed to SendChatMessage")
		return nil, model.ErrInternalServerError
	} else {
		message.ID = insertedID
	}

	s.EventsChan <- model.EventWrapper{
		Topic: consts.GlobalTopic,
		Event: model.NewChatMessageEvent(message).Event,
	}

	return &rpc.SendChatMessageResponse{
		MessageID: message.ID,
	}, nil
}
