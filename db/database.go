package db

import "abbysoft/gardarike-online/model"

type CharacterDatabase interface {
	BeginTransaction() error
	EndTransaction() error
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
	GetChatMessages(offset int, count int) ([]model.ChatMessage, error)
	GetMapChunk(x, y int64) (model.MapChunk, error)
	SaveOrUpdate(chunk model.MapChunk) error
}

type Database interface {
	CharacterDatabase
	AccountDatabase
	WorldDatabase
}
