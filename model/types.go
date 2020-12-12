package model

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"bytes"
	"encoding/gob"
	"fmt"
)

const GlobalTopic = "GLOBAL"
const MapChunkSize = 200

type EventWrapper struct {
	Topic string
	Event *rpc.Event
}

type Account struct {
	ID            int64  `db:"id"`
	Login         string `db:"login"`
	Password      string `db:"password"`
	Salt          string `db:"salt"`
	IsOnline      bool   `db:"is_online"`
	LastSessionID string `db:"last_session_id"`
}

type ChatMessage struct {
	ID     int64
	Sender string
	Text   string
}

func (c ChatMessage) ToRPC() *rpc.ChatMessage {
	return &rpc.ChatMessage{
		Id:     c.ID,
		Sender: c.Sender,
		Text:   c.Text,
	}
}

type Town struct {
	X          int64
	Y          int64
	OwnerName  string
	Population uint64
	Name       string
}

func (t Town) ToRPC() *rpc.Town {
	return &rpc.Town{
		X:          t.X,
		Y:          t.Y,
		Name:       t.Name,
		OwnerName:  t.OwnerName,
		Population: t.Population,
	}
}

type WorldMapChunk struct {
	X       int64
	Y       int64
	Width   int32
	Height  int32
	Data    []byte
	Towns   []Town
	Trees   uint64
	Stones  uint64
	Animals uint64
	Plants  uint64
}

func NewWorldMapChunkFromRPC(rpcChunk rpc.WorldMapChunk) WorldMapChunk {
	return WorldMapChunk{
		X:       rpcChunk.X,
		Y:       rpcChunk.Y,
		Width:   rpcChunk.Width,
		Height:  rpcChunk.Height,
		Data:    nil,
		Towns:   nil,
		Trees:   rpcChunk.Trees,
		Stones:  rpcChunk.Stones,
		Animals: rpcChunk.Animals,
		Plants:  rpcChunk.Plants,
	}
}

func (w WorldMapChunk) ToRPC() (*rpc.WorldMapChunk, error) {
	var terrain []float32
	decoder := gob.NewDecoder(bytes.NewBuffer(w.Data))
	if err := decoder.Decode(&terrain); err != nil {
		return nil, fmt.Errorf("failed to decode terrain data: %w", err)
	}

	mapChunk := &rpc.WorldMapChunk{
		X:       w.X,
		Y:       w.Y,
		Width:   w.Width,
		Height:  w.Height,
		Data:    terrain,
		Towns:   nil,
		Trees:   w.Trees,
		Stones:  w.Stones,
		Animals: w.Animals,
		Plants:  w.Plants,
	}

	for _, town := range w.Towns {
		mapChunk.Towns = append(mapChunk.Towns, town.ToRPC())
	}

	return mapChunk, nil
}

type Character struct {
	ID                int64
	AccountID         int64
	Name              string
	MaxPopulation     uint64 `db:"max_population"`
	CurrentPopulation uint64 `db:"current_population"`
	Towns             []Town
}

func (c Character) ToRPC() *rpc.Character {
	return &rpc.Character{
		Id:                c.ID,
		Name:              c.Name,
		MaxPopulation:     c.MaxPopulation,
		CurrentPopulation: c.CurrentPopulation,
	}
}

func NewChatMessageEvent(message rpc.ChatMessage) *rpc.Event {
	return &rpc.Event{
		Payload: &rpc.Event_ChatMessageEvent{
			ChatMessageEvent: &rpc.NewChatMessageEvent{
				Message: &message,
			},
		},
	}
}
