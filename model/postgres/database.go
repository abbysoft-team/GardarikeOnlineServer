package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"projectx-server/model"
)

type Database struct {
	db *sqlx.DB
}

func (d *Database) AddBuildingLocation(buildingLoc model.BuildingLocation) error {
	_, err := d.db.NamedExec(
		`INSERT INTO buildinglocations (building_id, owner_id, location)
VALUES (:building_id, :owner_id, :location)`, buildingLoc)
	return err
}

func (d *Database) GetBuildingOnLocation(location [3]float32) (result model.Building, err error) {
	err = d.db.Get(&result, "SELECT * FROM buildinglocations WHERE location=$1", location)
	return
}

func (d *Database) GetBuildings() (result []model.Building, err error) {
	err = d.db.Select(&result, "SELECT * FROM buildings")
	return
}

func (d *Database) GetBuildingLocations() (result []model.BuildingLocation, err error) {
	err = d.db.Select(&result, "SELECT * FROM buildinglocations")
	return
}

func (d *Database) GetCharacters(accountID int) (result []model.Character, err error) {
	err = d.db.Select(&result,
		`SELECT id, name, gold FROM accountcharacters as a
    INNER JOIN characters as c
        ON c.id = a.character_id
WHERE account_id = 1;`)

	return
}

func (d *Database) GetAccount(login string) (result model.Account, err error) {
	err = d.db.Get(&result, "SELECT * from accounts WHERE login = $1", login)
	return
}

func (d *Database) GetCharacter(id int) (result model.Character, err error) {
	err = d.db.Get(&result, "SELECT * FROM characters WHERE id = $1", id)
	return
}

func (d *Database) AddCharacter(character model.Character) error {
	_, err := d.db.NamedExec("INSERT INTO characters VALUES (DEFAULT, :name, :gold)", character)
	return err
}

func (d *Database) DeleteCharacter(id int) error {
	_, err := d.db.Exec("DELETE FROM characters WHERE id = $1", id)
	return err
}

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	EnableSSL bool
}

func NewDatabase(config Config) (model.Database, error) {
	var sslMode string
	if config.EnableSSL {
		sslMode = "verify-full"
	} else {
		sslMode = "disable"
	}

	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s",
			config.DBName,
			config.User,
			config.Password,
			config.Host,
			config.Port,
			sslMode))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return &Database{
		db: db,
	}, nil
}
