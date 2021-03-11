package logic

import (
	"abbysoft/gardarike-online/model"
	"abbysoft/gardarike-online/model/consts"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimpleLogic_SelectCharacter(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.SelectCharacterRequest{
		SessionID:   "sessionID",
		CharacterID: 2,
	}

	session.AccountID = 1
	character := model.Character{
		ID:                2,
		AccountID:         1,
		Name:              "test2",
		MaxPopulation:     1,
		CurrentPopulation: 1,
		Towns:             nil,
		Resources:         model.ResourcesPlaceTown,
	}

	towns := []model.Town{{ID: 1, X: 1, Y: 5, OwnerName: "test2", Name: "town", Population: 100}}

	db.On("GetCharacter", int64(2)).Return(character, nil)
	db.On("GetTowns", character.Name).Return(towns, nil)

	resp, err := logic.SelectCharacter(session, request)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 1, len(logic.EventsChan))

	event := <-logic.EventsChan
	require.Equal(t, model.NewSystemChatMessageEvent(consts.MessageCharacterAuthorized(character.Name)), event)

	require.NotEmpty(t, resp.Towns)
	require.Equal(t, model.ResourcesPlaceTown.ToRPC(), resp.Resources)

	db.AssertExpectations(t)
}

func TestSimpleLogic_SelectCharacter_WrongAccount(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.SelectCharacterRequest{
		SessionID:   "sessionID",
		CharacterID: 10,
	}

	session.AccountID = 1
	character := model.Character{
		ID:                10,
		AccountID:         2,
		Name:              "test10",
		MaxPopulation:     1,
		CurrentPopulation: 1,
		Towns:             nil,
	}

	db.On("GetCharacter", int64(10)).Return(character, nil)

	resp, err := logic.SelectCharacter(session, request)
	db.AssertExpectations(t)

	assert.EqualError(t, err, model.ErrForbidden.Error())
	assert.Nil(t, resp)
}
