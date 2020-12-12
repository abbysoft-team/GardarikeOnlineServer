package logic

import (
	"abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"errors"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	rand.Seed(time.Now().Unix())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *SimpleLogic) CreateAccount(session *PlayerSession, request *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, model.Error) {
	s.log.WithField("login", request.Login).Info("CreateAccount")

	salt := randStringBytes(10)
	id, err := s.db.AddAccount(request.Login, request.Password, salt)

	if err != nil && errors.Is(err, db.ErrDuplicatedUniqueKey) {
		return nil, model.ErrUsernameIsTaken
	} else if err != nil {
		s.log.WithError(err).Error("Failed to add new account to the db")
		return nil, model.ErrInternalServerError
	}

	return &rpc.CreateAccountResponse{
		Id: int32(id),
	}, nil
}
