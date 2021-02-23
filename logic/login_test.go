package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var invalidAccError = "invalid username/password combination"

func TestSimpleLogic_Login_InvalidUsername(t *testing.T) {
	logic, db, _ := NewLogicMock()
	request := &rpc.LoginRequest{
		Username: "John",
		Password: "hello",
	}

	db.On("GetAccount", "John").Return(model.Account{}, sql.ErrNoRows)

	_, err := logic.Login(request)

	assert.EqualError(t, err, model.ErrInvalidUserPassword.Error())
	db.AssertExpectations(t)
}

func TestSimpleLogic_Login_InvalidPassword(t *testing.T) {
	logic, db, _ := NewLogicMock()
	request := &rpc.LoginRequest{
		Username: "test",
		Password: "hello1",
	}

	db.On("GetAccount", "test").Return(model.Account{
		ID:            1,
		Login:         "test",
		Password:      saltPassword("hello", "salt"),
		Salt:          "salt",
		IsOnline:      true,
		LastSessionID: "",
	}, nil)

	_, err := logic.Login(request)
	assert.EqualError(t, err, model.ErrInvalidUserPassword.Error())
	db.AssertExpectations(t)
}

func TestSimpleLogic_Login(t *testing.T) {
	logic, db, _ := NewLogicMock()
	request := &rpc.LoginRequest{
		Username: "test",
		Password: "hello",
	}

	account := model.Account{
		ID:            1,
		Login:         "test",
		Password:      saltPassword("hello", "salt"),
		Salt:          "salt",
		IsOnline:      true,
		LastSessionID: "",
	}

	characters := []model.Character{
		{ID: 1, AccountID: account.ID, Name: "char"},
	}

	db.On("GetAccount", "test").Return(account, nil)
	db.On("GetCharacters", account.ID).Return(characters, nil)

	resp, err := logic.Login(request)
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, resp) {
		return
	}

	assert.NotEmpty(t, resp.Characters)
	db.AssertExpectations(t)
}
