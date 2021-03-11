package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) CreateCharacter(session *PlayerSession, request *rpc.CreateCharacterRequest) (*rpc.CreateCharacterResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
		"name":      request.Name,
	}).Info("CreateCharacter")

	tx := session.Tx

	id, err := tx.AddCharacter(request.Name)
	if err != nil {
		s.log.WithError(err).Error("Failed to create character")
		return nil, model.ErrInternalServerError
	}

	if err = tx.AddAccountCharacter(id, int(session.AccountID)); err != nil {
		s.log.WithError(err).Error("Failed to add account character")
		return nil, model.ErrInternalServerError
	}

	return &rpc.CreateCharacterResponse{
		Id: int64(id),
	}, nil
}
