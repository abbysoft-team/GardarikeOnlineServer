package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func (s *SimpleLogic) Login(request *rpc.LoginRequest) (*rpc.LoginResponse, model.Error) {
	s.log.WithField("login", request.GetUsername()).Info("Login request")

	acc, err := s.db.GetAccount(request.GetUsername())
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrInvalidUserPassword
	} else if err != nil {
		s.log.WithError(err).Error("Failed to get account from the database")
		return nil, model.ErrInternalServerError
	}

	hashedPass := md5.Sum([]byte(request.Password))
	saltedHash := fmt.Sprintf("%s%x%s", acc.Salt, string(hashedPass[:]), acc.Salt)
	finalPass := md5.Sum([]byte(saltedHash))

	if fmt.Sprintf("%x", string(finalPass[:])) != acc.Password {
		return nil, model.ErrInvalidUserPassword
	}

	chars, err := s.db.GetCharacters(acc.ID)
	if err != nil {
		s.log.WithError(err).WithField("accID", acc.ID).
			Error("Failed to get characters for account")
		return nil, model.ErrInternalServerError
	}

	session := NewPlayerSession()

	s.sessions[session.SessionID] = session

	s.log.WithFields(log.Fields{
		"accID":     acc.ID,
		"login":     acc.Login,
		"sessionID": session.SessionID,
	}).Info("User authorized on the server")

	var rpcChars []*rpc.Character
	for _, char := range chars {
		rpcChars = append(rpcChars, char.ToRPC())
	}

	return &rpc.LoginResponse{
		Characters: rpcChars,
		SessionID:  session.SessionID,
	}, nil
}
