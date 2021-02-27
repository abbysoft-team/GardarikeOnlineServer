package db

import (
	"abbysoft/gardarike-online/model"
)

type CharacterDatabase interface {
	GetCharacter(id int64) (model.Character, error)
	AddCharacter(name string, commit bool) (id int, err error)
	AddAccountCharacter(characterID, accountID int, commit bool) error
	DeleteCharacter(id int64, commit bool) error
	GetCharacters(accountID int64) ([]model.Character, error)
	UpdateCharacter(character model.Character, commit bool) error
}

type AccountDatabase interface {
	GetAccount(login string) (model.Account, error)
	AddAccount(login string, password string, salt string) (int, error)
}

type WorldDatabase interface {
	AddChatMessage(message model.ChatMessage) (int64, error)
	GetChatMessages(offset int, count int) ([]model.ChatMessage, error)
	GetMapChunk(x, y int64) (model.WorldMapChunk, error)
	SaveMapChunkOrUpdate(chunk model.WorldMapChunk, commit bool) error
	GetTowns(ownerName string) ([]model.Town, error)
	GetAllTowns() ([]model.Town, error)
	GetTownsForRect(xStart, xEnd, yStart, yEnd int) ([]model.Town, error)
	AddResourcesOrUpdate(resources model.Resources, commit bool) error
	GetResources(characterID int64) (model.Resources, error)
	AddTown(town model.Town, commit bool) error
}

type Database interface {
	CharacterDatabase
	AccountDatabase
	WorldDatabase
}
