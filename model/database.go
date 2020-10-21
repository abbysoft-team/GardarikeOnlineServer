package model

type CharacterDatabase interface {
	GetCharacter(id int) (Character, error)
	AddCharacter(character Character) error
	DeleteCharacter(id int) error
	GetCharacters(accountID int) ([]Character, error)
}

type AccountDatabase interface {
	GetAccount(login string) (Account, error)
}

type WorldDatabase interface {
	GetBuildingLocations() ([]BuildingLocation, error)
	GetBuildings() ([]Building, error)
	GetBuildingOnLocation(location [3]float32) (Building, error)
	AddBuildingLocation(buildingLoc BuildingLocation) error
}

type Database interface {
	CharacterDatabase
	AccountDatabase
	WorldDatabase
}
