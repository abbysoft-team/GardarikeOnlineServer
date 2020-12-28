package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type DatabaseMock struct {
	mock.Mock
}

func NewLogicMock() (*SimpleLogic, *DatabaseMock, *PlayerSession) {
	var s SimpleLogic
	var db DatabaseMock
	s.db = &db
	s.sessions = make(map[string]*PlayerSession)

	s.log = log.WithField("module", "test")

	session := NewPlayerSession(1)
	s.sessions[session.SessionID] = session

	return &s, &db, session
}

func (d *DatabaseMock) GetCharacter(id int64) (model.Character, error) {
	panic("implement me")
}

func (d *DatabaseMock) AddCharacter(character model.Character, commit bool) error {
	panic("implement me")
}

func (d *DatabaseMock) DeleteCharacter(id int64, commit bool) error {
	panic("implement me")
}

func (d *DatabaseMock) GetCharacters(accountID int64) ([]model.Character, error) {
	panic("implement me")
}

func (d *DatabaseMock) UpdateCharacter(character model.Character, commit bool) error {
	panic("implement me")
}

func (d *DatabaseMock) GetAccount(login string) (model.Account, error) {
	panic("implement me")
}

func (d *DatabaseMock) AddAccount(login string, password string, salt string) (int, error) {
	args := d.Called(login, password, salt)
	return args.Int(0), args.Error(1)
}

func (d *DatabaseMock) AddChatMessage(message model.ChatMessage) (int64, error) {
	panic("implement me")
}

func (d *DatabaseMock) GetChatMessages(offset int, count int) ([]model.ChatMessage, error) {
	panic("implement me")
}

func (d *DatabaseMock) GetMapChunk(x, y int64) (model.WorldMapChunk, error) {
	panic("implement me")
}

func (d *DatabaseMock) SaveOrUpdate(chunk model.WorldMapChunk, commit bool) error {
	panic("implement me")
}

func (d *DatabaseMock) GetTowns(ownerName string) ([]model.Town, error) {
	panic("implement me")
}

func TestSimpleLogic_CreateAccount(t *testing.T) {
	logic, db, session := NewLogicMock()
	request := &rpc.CreateAccountRequest{
		Login:    "login",
		Password: "password",
	}

	db.On("AddAccount", "login", mock.Anything, mock.Anything).
		Once().
		Return(1, nil)

	resp, err := logic.CreateAccount(session, request)
	if !assert.NoError(t, err, "request error is not nil") {
		return
	}

	if !assert.NotNil(t, resp, "response is nil") {
		return
	}

	assert.NotZero(t, resp.Id, "account ID is zero")
}
