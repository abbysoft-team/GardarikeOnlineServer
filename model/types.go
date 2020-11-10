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
	ID            int    `db:"id"`
	Login         string `db:"login"`
	Password      string `db:"password"`
	Salt          string `db:"salt"`
	IsOnline      bool   `db:"is_online"`
	LastSessionID string `db:"last_session_id"`
}

type Character struct {
	ID                int    `db:"id"`
	Name              string `db:"name"`
	Gold              uint64 `db:"gold"`
	MaxPopulation     uint64 `db:"max_population"`
	CurrentPopulation uint64 `db:"current_population"`
}

type Building struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	Cost            uint64 `db:"cost"`
	PopulationBonus uint64 `db:"population_bonus"`
}

type BuildingLocation struct {
	BuildingID int        `db:"building_id"`
	OwnerID    int        `db:"owner_id"`
	Location   [3]float32 `db:"location"`
}

type ChatMessage struct {
	MessageID  int    `db:"message_id"`
	SenderName string `db:"sender_name"`
	Text       string `db:"text"`
}

func (c *Character) ToRPC() *rpc.Character {
	return &rpc.Character{
		Id:                int32(c.ID),
		Name:              c.Name,
		Gold:              c.Gold,
		MaxPopulation:     c.MaxPopulation,
		CurrentPopulation: c.CurrentPopulation,
	}
}

func (building BuildingLocation) ToRPC() *rpc.Building {
	return &rpc.Building{
		Id:      int64(building.BuildingID),
		OwnerID: int64(building.OwnerID),
		Location: &rpc.Vector3D{
			X: building.Location[0],
			Y: building.Location[1],
			Z: building.Location[2],
		},
	}
}

func (message ChatMessage) ToRPC() *rpc.ChatMessage {
	return &rpc.ChatMessage{
		Id:     int64(message.MessageID),
		Sender: message.SenderName,
		Text:   message.Text,
	}
}

func NewPlaceBuildingEvent(buildingID int, ownerID int, location *rpc.Vector3D) *rpc.Event {
	return &rpc.Event{
		Payload: &rpc.Event_BuildingPlacedEvent{
			BuildingPlacedEvent: &rpc.BuildingPlacedEvent{
				BuildingID: int64(buildingID),
				OwnerID:    int64(ownerID),
				Location:   location,
			},
		},
	}
}

func NewChatMessageEvent(message ChatMessage) *rpc.Event {
	return &rpc.Event{
		Payload: &rpc.Event_ChatMessageEvent{
			ChatMessageEvent: &rpc.NewChatMessageEvent{
				Message: message.ToRPC(),
			},
		},
	}
}

func NewCharacterUpdatedEvent(char *Character) *rpc.Event {
	return &rpc.Event{
		Payload: &rpc.Event_CharacterUpdatedEvent{
			CharacterUpdatedEvent: &rpc.CharacterUpdatedEvent{
				NewState: char.ToRPC(),
			},
		},
	}
}

type MapChunk struct {
	X          int64  `db:"x"`
	Y          int64  `db:"y"`
	Terrain    []byte `db:"data"`
	TreesCount int64  `db:"trees_count"`
}

func (t MapChunk) ToRPC() (*rpc.Map, error) {
	var terrain []float32
	decoder := gob.NewDecoder(bytes.NewBuffer(t.Terrain))
	if err := decoder.Decode(&terrain); err != nil {
		return nil, fmt.Errorf("failed to decode terrain data: %w", err)
	}

	return &rpc.Map{
		Width:      MapChunkSize,
		Height:     MapChunkSize,
		Points:     terrain,
		Buildings:  nil,
		TreesCount: t.TreesCount,
	}, nil
}

func NewMapChunkFrom(rpcMap rpc.Map) (MapChunk, error) {
	var terrain []byte
	result := MapChunk{
		X:          0,
		Y:          0,
		TreesCount: rpcMap.TreesCount,
	}

	buffer := bytes.NewBuffer(terrain)
	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(&rpcMap.Points); err != nil {
		return MapChunk{}, fmt.Errorf("failed to encode map chunk: %w", err)
	}

	result.Terrain = buffer.Bytes()
	return result, nil
}
