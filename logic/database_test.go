package logic

import (
	"abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func NewLogicMock() (*SimpleLogic, *DatabaseTransactionMock, *PlayerSession) {
	var s SimpleLogic
	var db DatabaseMock
	s.db = &db
	s.sessions = make(map[string]*PlayerSession)

	s.log = log.WithField("module", "test")
	s.EventsChan = make(chan model.EventWrapper, 1)

	session := NewPlayerSession(1)
	s.sessions[session.SessionID] = session

	session.Tx = &db.DatabaseTransactionMock
	return &s, &db.DatabaseTransactionMock, session
}

func NewLogicMockWithTerrainGenerator() (*SimpleLogic, *DatabaseTransactionMock, *PlayerSession, *TerrainGeneratorMock) {
	logic, db, session := NewLogicMock()
	terrainGenerator := &TerrainGeneratorMock{}
	logic.generator = terrainGenerator

	return logic, db, session, terrainGenerator
}

type TerrainGeneratorMock struct {
	mock.Mock
}

func (t *TerrainGeneratorMock) GenerateTerrain(width int, height int, offsetX, offsetY float64) []float32 {
	args := t.Called(width, height, offsetX, offsetY)
	return args.Get(0).([]float32)
}

func (t *TerrainGeneratorMock) SetSeed(seed int64) {
}

type DatabaseTransactionMock struct {
	mock.Mock
	isCompleted bool
}

func (d *DatabaseTransactionMock) GetProductionRates(characterID int64) (model.Resources, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) AddOrUpdateResources(resources model.Resources) error {
	args := d.Called(resources)
	return args.Error(0)
}

func (d *DatabaseTransactionMock) AddOrUpdateProductionRates(rates model.Resources) error {
	panic("implement me")
}

func (d *DatabaseTransactionMock) AddTownBuilding(townID int64, building model.Building) error {
	panic("implement me")
}

func (d *DatabaseTransactionMock) GetAllBuildings() (map[int64]model.CharacterBuildings, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) EndTransaction() error {
	d.isCompleted = true
	return nil
}

func (d *DatabaseTransactionMock) IsCompleted() bool {
	return d.isCompleted
}

func (d *DatabaseTransactionMock) IsFailed() bool {
	return false
}

func (d *DatabaseTransactionMock) IsSucceed() bool {
	return true
}

func (d *DatabaseTransactionMock) SetAutoCommit(value bool) {

}

func (d *DatabaseTransactionMock) SetAutoRollBack(value bool) {

}

type DatabaseMock struct {
	DatabaseTransactionMock
}

func (d *DatabaseMock) BeginTransaction(autoCommit bool, autoRollBack bool) (db.DatabaseTransaction, error) {
	return d, nil
}

func (d *DatabaseTransactionMock) GetChunkRange() (model.ChunkRange, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) IncrementMapResources(resources model.ChunkResources) error {
	panic("implement me")
}

func (d *DatabaseTransactionMock) GetTownsForRect(xStart, xEnd, yStart, yEnd int) ([]model.Town, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) GetAllTowns() ([]model.Town, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) AddTown(town model.Town) error {
	args := d.Called(town)
	return args.Error(0)
}

func (d *DatabaseTransactionMock) AddResourcesOrUpdate(characterID int64, resources model.Resources) error {
	args := d.Called(resources)
	return args.Error(0)
}

func (d *DatabaseTransactionMock) GetResources(characterID int64) (model.Resources, error) {
	args := d.Called(characterID)
	return args.Get(0).(model.Resources), args.Error(1)
}

func (d *DatabaseTransactionMock) AddAccountCharacter(characterID, accountID int) error {
	args := d.Called(characterID, accountID)
	return args.Error(0)
}

func (d *DatabaseTransactionMock) GetCharacter(id int64) (model.Character, error) {
	args := d.Called(id)
	return args.Get(0).(model.Character), args.Error(1)
}

func (d *DatabaseTransactionMock) AddCharacter(name string) (int, error) {
	args := d.Called(name)
	return args.Int(0), args.Error(1)
}

func (d *DatabaseTransactionMock) DeleteCharacter(id int64) error {
	panic("implement me")
}

func (d *DatabaseTransactionMock) GetCharacters(accountID int64) ([]model.Character, error) {
	args := d.Called(accountID)
	return args.Get(0).([]model.Character), args.Error(1)
}

func (d *DatabaseTransactionMock) UpdateCharacter(character model.Character) error {
	args := d.Called(character)
	return args.Error(0)
}

func (d *DatabaseTransactionMock) GetAccount(login string) (model.Account, error) {
	args := d.Called(login)
	return args.Get(0).(model.Account), args.Error(1)
}

func (d *DatabaseTransactionMock) AddAccount(login string, password string, salt string) (int, error) {
	args := d.Called(login, password, salt)
	return args.Int(0), args.Error(1)
}

func (d *DatabaseTransactionMock) AddChatMessage(message model.ChatMessage) (int64, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) GetChatMessages(offset int, count int) ([]model.ChatMessage, error) {
	panic("implement me")
}

func (d *DatabaseTransactionMock) GetMapChunk(x, y int64) (model.WorldMapChunk, error) {
	args := d.Called(x, y)
	return args.Get(0).(model.WorldMapChunk), args.Error(1)
}

func (d *DatabaseTransactionMock) SaveMapChunkOrUpdate(chunk model.WorldMapChunk) error {
	args := d.Called(chunk)
	return args.Error(0)
}

func (d *DatabaseTransactionMock) GetTowns(ownerName string) ([]model.Town, error) {
	args := d.Called(ownerName)
	return args.Get(0).([]model.Town), args.Error(1)
}
