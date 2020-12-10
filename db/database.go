package db

import (
	"abbysoft/gardarike-online/model"
)

type CharacterDatabase interface {
	GetCharacter(id int) (model.Character, error)
	AddCharacter(character model.Character, commit bool) error
	DeleteCharacter(id int, commit bool) error
	GetCharacters(accountID int) ([]model.Character, error)
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
	SaveOrUpdate(chunk model.WorldMapChunk, commit bool) error
}

type Database interface {
	CharacterDatabase
	AccountDatabase
	WorldDatabase
}
