package game

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"projectx-server/model"
	rpc "projectx-server/rpc/generated"
)

var ErrInternalServerError = errors.New("internal server error")
var ErrInvalidUserPassword = errors.New("invalid username/password combination")

func (s *SimpleLogic) Login(request *rpc.LoginRequest) (*rpc.LoginResponse, error) {
	s.log.WithField("login", request.GetUsername()).Debug("Login request")

	acc, err := s.db.GetAccount(request.GetUsername())
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidUserPassword
	} else if err != nil {
		return nil, ErrInternalServerError
	}

	hashedPass := md5.Sum([]byte(request.Password))
	saltedHash := fmt.Sprintf("%s%x%s", acc.Salt, string(hashedPass[:]), acc.Salt)
	finalPass := md5.Sum([]byte(saltedHash))

	if fmt.Sprintf("%x", string(finalPass[:])) != acc.Password {
		return nil, ErrInvalidUserPassword
	}

	chars, err := s.db.GetCharacters(acc.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account characters: %w", err)
	}

	var rpcChars []*rpc.Character
	for _, char := range chars {
		rpcChars = append(rpcChars, model.ToRPCCharacter(char))
	}

	session := NewPlayerSession()
	s.sessions[session.SessionID] = session

	s.log.WithFields(log.Fields{
		"accID":     acc.ID,
		"login":     acc.Login,
		"sessionID": session.SessionID,
	}).Info("User authorized on the server")

	return &rpc.LoginResponse{
		Characters: rpcChars,
		SessionID:  session.SessionID,
	}, nil
}
