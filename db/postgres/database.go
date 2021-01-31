package postgres

import (
	"abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	"fmt"
	"github.com/jmoiron/sqlx"

	pq "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func (d *Database) beginTransaction() error {
	tx, err := d.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	d.tx = tx
	return nil
}

func (d *Database) endTransaction() error {
	if d.tx != nil {
		if err := d.tx.Commit(); err != nil {
			return fmt.Errorf("failed to end transaction: %w", err)
		}

		d.tx = nil
	}

	return nil
}

type transactionFunc func(t *sqlx.Tx) error

func (d *Database) AddTown(town model.Town, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := t.NamedExec(
			`INSERT INTO towns VALUES (:x, :y, :owner_name, :population, :name)`, town)
		return err
	}, commit)
}

func (d *Database) AddResourcesOrUpdate(resources model.Resources, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := t.NamedExec(
			`INSERT INTO resources VALUES (:character_id, DEFAULT, DEFAULT, DEFAULT, DEFAULT) 
ON CONFLICT (character_id) DO UPDATE SET
stone = :stone, food = :food, leather = :leather, wood = :wood
`, resources)
		return err
	}, commit)
}

func (d *Database) GetResources(characterID int64) (result model.Resources, err error) {
	err = d.db.Get(&result, "SELECT * FROM resources WHERE character_id=$1", characterID)
	return
}

func (d *Database) AddAccountCharacter(characterID, accountID int, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := t.Exec("INSERT INTO account_characters VALUES ($1, $2)", accountID, characterID)
		return err
	}, commit)
}

func (d *Database) AddAccount(login string, password string, salt string) (id int, err error) {
	err = d.db.Get(&id,
		"INSERT INTO accounts VALUES (DEFAULT, $1, $2, $3, DEFAULT, DEFAULT) RETURNING id", login, password, salt)
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" { // Unique key constraint
			return 0, db.ErrDuplicatedUniqueKey
		}
	}
	return
}

func (d *Database) WithTransaction(function transactionFunc, commit bool) error {
	if d.tx == nil {
		if err := d.beginTransaction(); err != nil {
			return err
		}
	}

	if err := function(d.tx); err != nil {
		if rollbackErr := d.tx.Rollback(); rollbackErr != nil {
			d.tx = nil
			return fmt.Errorf("%w: (and failed to rollback: %v)", err, rollbackErr)
		}

		d.tx = nil
		return err
	}

	if commit {
		if err := d.endTransaction(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) GetAllTowns() (result []model.Town, err error) {
	err = d.db.Select(&result, "SELECT * FROM towns")
	return
}

func (d *Database) GetTowns(ownerName string) (result []model.Town, err error) {
	err = d.db.Select(&result, "SELECT * FROM towns WHERE owner_name=$1", ownerName)
	return
}

func (d *Database) SaveOrUpdate(chunk model.WorldMapChunk, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := t.NamedQuery(
			`INSERT INTO chunks (x, y, data, trees, stones, animals, plants) VALUES 
                                      (:x, :y, :data, :trees, :stones, :animals, :plants)
			   ON CONFLICT (x, y) DO UPDATE 
			   SET trees = :trees,
			   stones = :stones,
			   animals = :animals,
			   plants = :plants`, chunk)

		return err
	}, commit)
}

func (d *Database) GetMapChunk(x, y int64) (result model.WorldMapChunk, err error) {
	err = d.db.Get(&result, "SELECT * FROM chunks WHERE x=$1 AND y=$2", x, y)

	return
}

func (d *Database) GetChatMessages(offset int, count int) (result []model.ChatMessage, err error) {
	err = d.db.Select(&result,
		"SELECT * FROM chat_messages ORDER BY message_id DESC OFFSET $1 LIMIT $2",
		offset, count)
	return
}

func (d *Database) AddChatMessage(message model.ChatMessage) (id int64, err error) {
	err = d.db.Get(&id,
		"INSERT INTO chat_messages (sender_name, text) VALUES ($1, $2) RETURNING message_id",
		message.Sender, message.Text)
	return
}

func (d *Database) UpdateCharacter(character model.Character, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := d.db.NamedExec(
			`UPDATE characters SET 
                      name=:name, max_population=:maxPopulation, current_population=:currentPopulation
			   WHERE id=:Id`, &character)
		return err
	}, commit)
}

func (d *Database) GetCharacters(accountID int64) (result []model.Character, err error) {
	err = d.db.Select(&result, `SELECT c.* FROM account_characters as a
    INNER JOIN characters as c
        ON c.id = a.character_id
WHERE account_id = $1`, accountID)

	return
}

func (d *Database) GetAccount(login string) (result model.Account, err error) {
	err = d.db.Get(&result, "SELECT * from accounts WHERE login = $1", login)
	return
}

func (d *Database) GetCharacter(id int64) (result model.Character, err error) {
	err = d.db.Get(&result,
		`SELECT c.*, ac.account_id FROM characters c 
    JOIN account_characters ac on c.id = ac.character_id WHERE c.id = $1`, id)
	return
}

func (d *Database) AddCharacter(name string, commit bool) (id int, err error) {
	err = d.WithTransaction(func(t *sqlx.Tx) error {
		err = t.Get(&id, "INSERT INTO characters VALUES (DEFAULT, $1, DEFAULT, DEFAULT) RETURNING id", name)
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" { // Unique key constraint
					return db.ErrDuplicatedUniqueKey
				}
			}
		}
		return err
	}, commit)

	return
}

func (d *Database) DeleteCharacter(id int64, commit bool) error {
	return d.WithTransaction(func(t *sqlx.Tx) error {
		_, err := d.db.Exec("DELETE FROM characters WHERE id = $1", id)
		return err
	}, commit)
}

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	EnableSSL bool
}

func NewDatabase(config Config) (db.Database, error) {
	var sslMode string
	if config.EnableSSL {
		sslMode = "verify-full"
	} else {
		sslMode = "disable"
	}

	database, err := sqlx.Connect("postgres",
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
		db: database,
	}, nil
}
