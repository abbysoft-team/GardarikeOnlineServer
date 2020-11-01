package model

import (
	rpc "abbysoft/gardarike-online/rpc/generated"
)

type Account struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Salt     string `db:"salt"`
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
		Id:   int32(c.ID),
		Name: c.Name,
		Gold: c.Gold,
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
