package logic

import (
	"abbysoft/gardarike-online/model"
	log "github.com/sirupsen/logrus"
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

func (d *DatabaseMock) GetCharacter(id int64) (model.Character, error) {
	panic("implement me")
}

func (d *DatabaseMock) AddCharacter(name string) (int, error) {
	args := d.Called(name)
	return args.Int(0), args.Error(1)
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
