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

func saltPassword(password, salt string) string {
	hashedPass := md5.Sum([]byte(password))
	saltedHash := fmt.Sprintf("%s%x%s", salt, string(hashedPass[:]), salt)
	finalPass := md5.Sum([]byte(saltedHash))

	return fmt.Sprintf("%x", string(finalPass[:]))
}

func (s *SimpleLogic) Login(request *rpc.LoginRequest) (*rpc.LoginResponse, model.Error) {
	s.log.WithField("login", request.GetUsername()).Info("Login request")

	acc, err := s.db.GetAccount(request.GetUsername())
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, model.ErrInvalidUserPassword
	} else if err != nil {
		s.log.WithError(err).Error("Failed to get account from the database")
		return nil, model.ErrInternalServerError
	}

	if saltPassword(request.Password, acc.Salt) != acc.Password {
		return nil, model.ErrInvalidUserPassword
	}

	chars, err := s.db.GetCharacters(acc.ID)
	if err != nil {
		s.log.WithError(err).WithField("accID", acc.ID).
			Error("Failed to get characters for account")
		return nil, model.ErrInternalServerError
	}

	session := NewPlayerSession(acc.ID)

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
