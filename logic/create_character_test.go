package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleLogic_CreateCharacter(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.CreateCharacterRequest{
		Name: "test",
	}

	session.AccountID = 2

	db.On("AddCharacter", "test").Return(1, nil)
	db.On("AddAccountCharacter", 1, 2).Return(nil)

	resp, err := logic.CreateCharacter(session, request)

	db.AssertExpectations(t)

	assert.NoError(t, err, "failed to create character")
	if !assert.NotNil(t, resp) {
		return
	}

	assert.NotZero(t, resp.Id)
}

func TestSimpleLogic_CreateCharacter_DuplicateLogin(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.CreateCharacterRequest{
		Name: "test",
	}

	db.On("AddCharacter", "test").
		Return(0, errors.New("some error"))

	resp, err := logic.CreateCharacter(session, request)
	db.AssertExpectations(t)

	assert.EqualError(t, err, model.ErrInternalServerError.Error(), "error isn't internal error")
	assert.Nil(t, resp)
}
