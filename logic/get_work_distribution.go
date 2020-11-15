package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) GetWorkDistribution(session *PlayerSession, request *rpc.GetWorkDistributionRequest) (*rpc.GetWorkDistributionResponse, model.Error) {
	s.log.WithFields(log.Fields{
		"sessionID": session.SessionID,
	}).Info("GetWorkDistribution")

	return &session.WorkDistribution, nil
}
