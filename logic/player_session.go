package logic

import (
	"abbysoft/gardarike-online/model"
	"github.com/google/uuid"
	"sync"
)

type PlayerSession struct {
	SessionID         string
	SelectedCharacter *model.Character
	Mutex             sync.Mutex
}

func NewPlayerSession() *PlayerSession {
	return &PlayerSession{
		SessionID:         uuid.New().String(),
		SelectedCharacter: nil,
	}
}
