package logic

import (
	"abbysoft/gardarike-online/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

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

type DatabaseMock struct {
	mock.Mock
}

func (d *DatabaseMock) GetResources(characterID int64) (model.Resources, error) {
	args := d.Called(characterID)
	return args.Get(0).(model.Resources), args.Error(1)
}

func (d *DatabaseMock) AddAccountCharacter(characterID, accountID int) error {
	args := d.Called(characterID, accountID)
	return args.Error(0)
}

func (d *DatabaseMock) GetCharacter(id int64) (model.Character, error) {
	args := d.Called(id)
	return args.Get(0).(model.Character), args.Error(1)
}

func (d *DatabaseMock) AddCharacter(name string) (int, error) {
	args := d.Called(name)
	return args.Int(0), args.Error(1)
}

func (d *DatabaseMock) DeleteCharacter(id int64, commit bool) error {
	panic("implement me")
}

func (d *DatabaseMock) GetCharacters(accountID int64) ([]model.Character, error) {
	args := d.Called(accountID)
	return args.Get(0).([]model.Character), args.Error(1)
}

func (d *DatabaseMock) UpdateCharacter(character model.Character, commit bool) error {
	panic("implement me")
}

func (d *DatabaseMock) GetAccount(login string) (model.Account, error) {
	args := d.Called(login)
	return args.Get(0).(model.Account), args.Error(1)
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
	args := d.Called(ownerName)
	return args.Get(0).([]model.Town), args.Error(1)
}
