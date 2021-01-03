package logic

import (
	db2 "abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type DatabaseMock struct {
	mock.Mock
}

func TestSimpleLogic_CreateAccount(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.CreateAccountRequest{
		Login:    "login",
		Password: "password",
	}

	db.On("AddAccount", "login", mock.Anything, mock.Anything).Once().
		Return(1, nil).Run(func(args mock.Arguments) {
		pass := args.String(1)
		salt := args.String(2)

		expectedPass := request.Password

		assert.Equal(t, saltPassword(expectedPass, salt), pass, "password not salted properly")
	})

	resp, err := logic.CreateAccount(session, request)
	if !assert.NoError(t, err, "request error is not nil") {
		return
	}

	if !assert.NotNil(t, resp, "response is nil") {
		return
	}

	assert.NotZero(t, resp.Id, "account ID is zero")
}

func TestSimpleLogic_CreateAccount_AlreadyExists(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.CreateAccountRequest{
		Login:    "login",
		Password: "password",
	}

	db.On("AddAccount", "login", mock.Anything, mock.Anything).Once().
		Return(0, db2.ErrDuplicatedUniqueKey)

	_, err := logic.CreateAccount(session, request)
	assert.EqualError(t, err, model.ErrUsernameIsTaken.Error())
}
