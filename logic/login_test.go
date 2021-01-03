package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"testing"
)

type databaseMock struct {
	getCharacterInvocations int
	getAccountInvocations   int
}

func (d *databaseMock) AddAccountCharacter(characterID, accountID int) error {
	panic("implement me")
}

func (d *databaseMock) GetTowns(ownerName string) ([]model.Town, error) {
	panic("implement me")
}

func (d *databaseMock) AddAccount(login string, password string, salt string) (int, error) {
	panic("implement me")
}

func (d *databaseMock) GetCharacter(id int64) (model.Character, error) {
	panic("implement me")
}

func (d *databaseMock) AddCharacter(name string) (int, error) {
	panic("implement me")
}

func (d *databaseMock) DeleteCharacter(id int64, commit bool) error {
	panic("implement me")
}

func (d *databaseMock) UpdateCharacter(character model.Character, commit bool) error {
	panic("implement me")
}

func (d *databaseMock) GetAccount(login string) (model.Account, error) {
	d.getAccountInvocations++

	if login != "test" {
		return model.Account{}, sql.ErrNoRows
	}

	return model.Account{
		ID:       1,
		Login:    "test",
		Password: "89cb1297f75457552c074c08e8e28a93",
		Salt:     "salt",
	}, nil
}

func (d *databaseMock) GetMapChunk(x, y int64) (model.WorldMapChunk, error) {
	panic("implement me")
}

func (d *databaseMock) SaveOrUpdate(chunk model.WorldMapChunk, commit bool) error {
	panic("implement me")
}

func (d *databaseMock) AddChatMessage(message model.ChatMessage) (int64, error) {
	panic("implement me")
}

func (d *databaseMock) GetChatMessages(offset int, count int) ([]model.ChatMessage, error) {
	panic("implement me")
}

func (d *databaseMock) GetCharactersInvocations() int {
	defer func() {
		d.getCharacterInvocations = 0
	}()

	return d.getCharacterInvocations
}

func (d *databaseMock) GetAccountInvocations() int {
	defer func() {
		d.getAccountInvocations = 0
	}()

	return d.getAccountInvocations
}

func (d *databaseMock) GetCharacters(accountID int64) ([]model.Character, error) {
	d.getCharacterInvocations++

	return []model.Character{
		{ID: 1, AccountID: 1, Name: "jack", MaxPopulation: 10, CurrentPopulation: 10},
		{2, 1, "lenny", 100, 10, nil},
		{3, 1, "michel", 100, 10, nil},
	}, nil
}

func newLogicMock() (*SimpleLogic, *databaseMock) {
	var s SimpleLogic
	var db databaseMock
	s.db = &db
	s.sessions = make(map[string]*PlayerSession)

	s.log = log.WithField("module", "test")

	return &s, &db
}

var invalidAccError = "invalid username/password combination"

func TestSimpleLogic_Login_InvalidUsername(t *testing.T) {
	logic, db := newLogicMock()
	request := &rpc.LoginRequest{
		Username: "John",
		Password: "hello",
	}

	_, err := logic.Login(request)
	if err == nil || err != model.ErrInvalidUserPassword {
		t.Fatalf("Login expected to return ErrInvalidUserPassword error but err is %v", err)
	}

	if invocs := db.GetAccountInvocations(); invocs != 1 {
		t.Fatalf("GetAccount invocations expected to be 1, but got %d", invocs)
	}
}

func TestSimpleLogic_Login_InvalidPassword(t *testing.T) {
	logic, db := newLogicMock()
	request := &rpc.LoginRequest{
		Username: "test",
		Password: "hello1",
	}

	_, err := logic.Login(request)
	if err == nil || err != model.ErrInvalidUserPassword {
		t.Fatalf("Login expected to return ErrInvalidUserPassword error but err is %v", err)
	}

	if invocs := db.GetAccountInvocations(); invocs != 1 {
		t.Fatalf("GetAccount invocations expected to be 1, but got %d", invocs)
	}
}

func TestSimpleLogic_Login(t *testing.T) {
	logic, db := newLogicMock()
	request := &rpc.LoginRequest{
		Username: "test",
		Password: "hello",
	}

	_, err := logic.Login(request)
	if err != nil {
		t.Fatalf("err expected to be nil but %v found", err)
	}

	if invocs := db.GetAccountInvocations(); invocs != 1 {
		t.Fatalf("GetAccount invocations expected to be 1, but got %d", invocs)
	}

	if invocs := db.GetCharactersInvocations(); invocs != 1 {
		t.Fatalf("GetCharacters invocations expected to be 1, but got %d", invocs)
	}
}
