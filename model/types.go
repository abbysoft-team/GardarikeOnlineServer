package model

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
	"bytes"
	"encoding/gob"
	"fmt"
)

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
	ID       int64
	Sender   string
	Text     string
	IsSystem bool `db:"is_system"`
}

func (c ChatMessage) ToRPC() *rpc.ChatMessage {
	var messageType rpc.ChatMessage_Type
	if c.IsSystem {
		messageType = rpc.ChatMessage_SYSTEM
	} else {
		messageType = rpc.ChatMessage_NORMAL
	}

	return &rpc.ChatMessage{
		Id:     c.ID,
		Sender: c.Sender,
		Text:   c.Text,
		Type:   messageType,
	}
}

type Town struct {
	X          int64
	Y          int64
	OwnerName  string `db:"owner_name"`
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

func NewWorldMapChunkFromRPC(rpcChunk rpc.WorldMapChunk) (WorldMapChunk, error) {
	var terrain []byte
	result := WorldMapChunk{
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

	buffer := bytes.NewBuffer(terrain)
	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(&rpcChunk.Data); err != nil {
		return WorldMapChunk{}, fmt.Errorf("failed to encode map chunk: %w", err)
	}

	result.Data = buffer.Bytes()

	return result, nil
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
	AccountID         int64 `db:"account_id"`
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

type Resources struct {
	CharacterID int `db:"character_id"`
	Wood        uint64
	Food        uint64
	Stone       uint64
	Leather     uint64
}

func (r Resources) ToRPC() *rpc.Resources {
	return &rpc.Resources{
		Wood:    r.Wood,
		Stone:   r.Stone,
		Food:    r.Food,
		Leather: r.Leather,
	}
}
