package logic

import (
	"abbysoft/gardarike-online/model"
	"github.com/google/uuid"
	"sync"
	"time"
)

type PlayerSession struct {
	SessionID         string
	SelectedCharacter *model.Character
	Mutex             sync.Mutex
	LastRequestTime   time.Time
}

func NewPlayerSession() *PlayerSession {
	return &PlayerSession{
		SessionID:         uuid.New().String(),
		SelectedCharacter: nil,
		LastRequestTime:   time.Now(),
	}
}
