package model

import rpc "projectx-server/rpc/generated"

type Account struct {
	ID       int    `db:"id"`
	Login    string `db:"login"`
	Password string `db:"password"`
	Salt     string `db:"salt"`
}

type Character struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Gold uint64 `db:"gold"`
}

type Building struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Cost int    `db:"cost"`
}

type BuildingLocation struct {
	BuildingID int        `db:"building_id"`
	OwnerID    int        `db:"owner_id"`
	Location   [3]float32 `db:"location"`
}

func (char Character) ToRPC() *rpc.Character {
	return &rpc.Character{
		Id:   int32(char.ID),
		Name: char.Name,
		Gold: char.Gold,
	}
}

func (building BuildingLocation) ToRPC() *rpc.Building {
	return &rpc.Building{
		Id:      int64(building.BuildingID),
		OwnerID: int64(building.OwnerID),
		Location: &rpc.Vector3D{
			X: float32(building.Location[0]),
			Y: float32(building.Location[1]),
			Z: float32(building.Location[2]),
		},
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
