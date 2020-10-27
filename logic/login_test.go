package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

type databaseMock struct {
	getCharacterInvocations int
	getAccountInvocations   int
}

func (d *databaseMock) UpdateCharacter(character model.Character) error {
	panic("implement me")
}

func (d *databaseMock) GetBuildingLocations() ([]model.BuildingLocation, error) {
	panic("implement me")
}

func (d *databaseMock) GetBuildings() ([]model.Building, error) {
	panic("implement me")
}

func (d *databaseMock) GetBuildingLocation(location [3]float32) (model.BuildingLocation, error) {
	panic("implement me")
}

func (d *databaseMock) AddBuildingLocation(buildingLoc model.BuildingLocation) error {
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

func (d *databaseMock) GetCharacter(id int) (model.Character, error) {
	panic("implement me")
}

func (d *databaseMock) AddCharacter(character model.Character) error {
	panic("implement me")
}

func (d *databaseMock) DeleteCharacter(id int) error {
	panic("implement me")
}

func (d *databaseMock) GetCharacters(accountID int) ([]model.Character, error) {
	d.getCharacterInvocations++

	return []model.Character{
		{ID: 1, Name: "jack", Gold: 100},
		{2, "lenny", 100},
		{3, "michel", 100},
	}, nil
}

func (d *databaseMock) GetAccount(login string) (model.Account, error) {
	d.getAccountInvocations++

	if login != "test" {
		return model.Account{}, fmt.Errorf("not found")
	}

	return model.Account{
		ID:       1,
		Login:    "test",
		Password: "89cb1297f75457552c074c08e8e28a93",
		Salt:     "salt",
	}, nil
}

func newLogicMock() (*SimpleLogic, *databaseMock) {
	var s SimpleLogic
	var db databaseMock
	s.db = &db
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
	if err == nil || err.Error() != invalidAccError {
		t.Fatalf("Login expected to return invalid acc error but err is %v", err)
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
	if err == nil || err.Error() != invalidAccError {
		t.Fatalf("Login expected to return invalid acc error but err is %v", err)
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
