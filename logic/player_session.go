package logic

import (
	"abbysoft/gardarike-online/model"
	"github.com/google/uuid"
)

type PlayerSession struct {
	SessionID         string
	SelectedCharacter *model.Character
}

func NewPlayerSession() *PlayerSession {
	return &PlayerSession{
		SessionID:         uuid.New().String(),
		SelectedCharacter: nil,
	}
}
