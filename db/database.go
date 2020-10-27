package db

import "abbysoft/gardarike-online/model"

type CharacterDatabase interface {
	GetCharacter(id int) (model.Character, error)
	AddCharacter(character model.Character) error
	DeleteCharacter(id int) error
	GetCharacters(accountID int) ([]model.Character, error)
	UpdateCharacter(character model.Character) error
}

type AccountDatabase interface {
	GetAccount(login string) (model.Account, error)
}

type WorldDatabase interface {
	GetBuildingLocations() ([]model.BuildingLocation, error)
	GetBuildings() ([]model.Building, error)
	GetBuildingLocation(location [3]float32) (model.BuildingLocation, error)
	AddBuildingLocation(buildingLoc model.BuildingLocation) error
	AddChatMessage(message model.ChatMessage) (int64, error)
}

type Database interface {
	CharacterDatabase
	AccountDatabase
	WorldDatabase
}
