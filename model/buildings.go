package model

import rpc "abbysoft/gardarike-online/rpc/generated"

type Location2D struct {
	X float32
	Y float32
}

type Building struct {
	ID       rpc.BuildingType
	Name     string
	Cost     Resources
	Location Location2D
}

// CharacterBuildings - number of buildings of each type
type CharacterBuildings map[rpc.BuildingType]uint64

var (
	Buildings = map[rpc.BuildingType]Building{
		rpc.BuildingType_HOUSE: {
			ID:   rpc.BuildingType_HOUSE,
			Name: "house",
			Cost: Resources{Wood: 30, Food: 10, Stone: 15, Leather: 20},
		},
		rpc.BuildingType_QUARRY: {
			ID:   rpc.BuildingType_QUARRY,
			Name: "quarry",
			Cost: Resources{Wood: 100, Food: 50, Stone: 0, Leather: 80},
		},
	}
)

func IsValidBuildingType(typeValue int32) bool {
	_, found := rpc.BuildingType_name[typeValue]
	return found
}
