package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) GetResources(
	session *PlayerSession, request *rpc.GetResourcesRequest) (*rpc.GetResourcesResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
	}).Info("GetResources")

	if session.SelectedCharacter == nil {
		return nil, model.ErrCharacterNotSelected
	}

	resources, err := s.db.GetResources(session.SelectedCharacter.ID)
	if err != nil {
		s.log.WithError(err).Error("Failed to get resources")
		return nil, model.ErrInternalServerError
	}

	return &rpc.GetResourcesResponse{
		Resources: resources.ToRPC(),
	}, nil
}
