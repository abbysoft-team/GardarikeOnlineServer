package game

import (
	"github.com/google/uuid"
	"projectx-server/model"
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
