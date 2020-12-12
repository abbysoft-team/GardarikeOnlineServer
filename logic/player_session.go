package logic

import (
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
	"github.com/google/uuid"
	"sync"
	"time"
)

type PlayerSession struct {
	SessionID         string
	AccountID         int64
	SelectedCharacter *model.Character
	Mutex             sync.Mutex
	LastRequestTime   time.Time
	WorkDistribution  rpc.GetWorkDistributionResponse
}

func NewPlayerSession(accountID int64) *PlayerSession {
	return &PlayerSession{
		SessionID:         uuid.New().String(),
		SelectedCharacter: nil,
		LastRequestTime:   time.Now(),
		WorkDistribution: rpc.GetWorkDistributionResponse{
			IdleCount:       0,
			WoodcutterCount: 0,
		},
	}
}
