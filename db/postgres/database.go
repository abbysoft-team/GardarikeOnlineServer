package postgres

import (
	"abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	"fmt"
	"github.com/jmoiron/sqlx"
	pq "github.com/lib/pq"
)

type DatabaseTransaction struct {
	tx           *sqlx.Tx
	autoCommit   bool
	autoRollBack bool
	isRolledBack bool
	isCommitted  bool
}

func (d *DatabaseTransaction) SetAutoCommit(value bool) {
	d.autoCommit = value
}

func (d *DatabaseTransaction) SetAutoRollBack(value bool) {
	d.autoRollBack = value
}

func (d *DatabaseTransaction) IsCompleted() bool {
	return d.isRolledBack || d.isCommitted
}

func (d *DatabaseTransaction) IsFailed() bool {
	return d.isRolledBack
}

func (d *DatabaseTransaction) IsSucceed() bool {
	return d.isCommitted
}

type Database struct {
	db *sqlx.DB
}

func (d *Database) BeginTransaction(autoCommit, autoRollBack bool) (db.DatabaseTransaction, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &DatabaseTransaction{tx: tx, autoCommit: autoCommit, autoRollBack: autoRollBack}, nil
}

func (d *DatabaseTransaction) EndTransaction() error {
	if d.IsCompleted() {
		return nil
	}

	if d.tx != nil {
		if err := d.tx.Commit(); err != nil {
			return fmt.Errorf("failed to end transaction: %w", err)
		}
		d.tx = nil
		d.isCommitted = true
		return nil
	} else {
		return fmt.Errorf("transaction is not started")
	}
}

type transactionFunc func(t *sqlx.Tx) error

func (d *DatabaseTransaction) handleError(err error) error {
	if err == nil {
		if d.autoCommit {
			if commitErr := d.tx.Commit(); commitErr != nil {
				return fmt.Errorf("failed to commit changes: %w", err)
			} else {
				d.tx = nil
				d.isCommitted = true
			}
		}

		return err
	}

	if d.autoRollBack {
		if rollbackErr := d.tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction failed: %w (rollback also failed: %v)", err, rollbackErr)
		} else {
			d.tx = nil
			d.isRolledBack = true
		}
	}

	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" { // Unique key constraint
			return db.ErrDuplicatedUniqueKey
		}
	}

	return err
}

func (d *DatabaseTransaction) IncrementMapResources(resources model.ChunkResources) error {
	_, err := d.tx.NamedExec(
		"UPDATE chunks SET stones = stones + :stones, trees = trees + :trees, animals = animals + :animals, plants = plants + :plants",
		resources)

	return d.handleError(err)
}

func (d *DatabaseTransaction) GetChunkRange() (result model.ChunkRange, err error) {
	err = d.tx.Get(&result,
		"SELECT MAX(x) as max_x, MIN(y) as min_x, MAX(y) as max_y, MIN(y) as min_y FROM chunks")
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) GetTownsForRect(xStart, xEnd, yStart, yEnd int) (results []model.Town, err error) {
	err = d.tx.Select(&results,
		"SELECT * FROM towns WHERE (x BETWEEN $1 AND $2) AND (y BETWEEN $3 AND $4)",
		xStart, xEnd, yStart, yEnd)
	return results, d.handleError(err)
}

func (d *DatabaseTransaction) AddTown(town model.Town) error {
	_, err := d.tx.NamedExec(
		`INSERT INTO towns VALUES (:x, :y, :owner_name, :population, :name)`, town)
	return d.handleError(err)
}

func (d *DatabaseTransaction) AddResourcesOrUpdate(resources model.Resources) error {
	_, err := d.tx.NamedExec(
		`INSERT INTO resources VALUES (:character_id, DEFAULT, DEFAULT, DEFAULT, DEFAULT) 
ON CONFLICT (character_id) DO UPDATE SET
stone = :stone, food = :food, leather = :leather, wood = :wood
`, resources)
	return d.handleError(err)
}

func (d *DatabaseTransaction) GetResources(characterID int64) (result model.Resources, err error) {
	err = d.tx.Get(&result, "SELECT * FROM resources WHERE character_id=$1", characterID)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) AddAccountCharacter(characterID, accountID int) error {
	_, err := d.tx.Exec("INSERT INTO account_characters VALUES ($1, $2)", accountID, characterID)
	return d.handleError(err)
}

func (d *DatabaseTransaction) AddAccount(login string, password string, salt string) (id int, err error) {
	err = d.tx.Get(&id,
		"INSERT INTO accounts VALUES (DEFAULT, $1, $2, $3, DEFAULT, DEFAULT) RETURNING id", login, password, salt)

	return id, d.handleError(err)
}

func (d *DatabaseTransaction) GetAllTowns() (result []model.Town, err error) {
	err = d.tx.Select(&result, "SELECT * FROM towns")
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) GetTowns(ownerName string) (result []model.Town, err error) {
	err = d.tx.Select(&result, "SELECT * FROM towns WHERE owner_name=$1", ownerName)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) SaveMapChunkOrUpdate(chunk model.WorldMapChunk) error {
	_, err := d.tx.NamedQuery(
		`INSERT INTO chunks (x, y, data, trees, stones, animals, plants) VALUES 
                                      (:x, :y, :data, :trees, :stones, :animals, :plants)
			   ON CONFLICT (x, y) DO UPDATE 
			   SET trees = :trees,
			   stones = :stones,
			   animals = :animals,
			   plants = :plants`, chunk)

	return d.handleError(err)
}

func (d *DatabaseTransaction) GetMapChunk(x, y int64) (result model.WorldMapChunk, err error) {
	err = d.tx.Get(&result, "SELECT * FROM chunks WHERE x=$1 AND y=$2", x, y)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) GetChatMessages(offset int, count int) (result []model.ChatMessage, err error) {
	err = d.tx.Select(&result,
		"SELECT * FROM chat_messages ORDER BY message_id DESC OFFSET $1 LIMIT $2",
		offset, count)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) AddChatMessage(message model.ChatMessage) (id int64, err error) {
	err = d.tx.Get(&id,
		"INSERT INTO chat_messages (sender_name, text) VALUES ($1, $2) RETURNING message_id",
		message.Sender, message.Text)
	return id, d.handleError(err)
}

func (d *DatabaseTransaction) UpdateCharacter(character model.Character) error {
	_, err := d.tx.NamedExec(
		`UPDATE characters SET 
			  name=:name, 
			  max_population=:max_population, 
			  current_population=:current_population
         WHERE id=:id`, &character)
	return d.handleError(err)
}

func (d *DatabaseTransaction) GetCharacters(accountID int64) (result []model.Character, err error) {
	err = d.tx.Select(&result, `SELECT c.* FROM account_characters as a
    INNER JOIN characters as c
        ON c.id = a.character_id
WHERE account_id = $1`, accountID)

	return result, d.handleError(err)
}

func (d *DatabaseTransaction) GetAccount(login string) (result model.Account, err error) {
	err = d.tx.Get(&result, "SELECT * from accounts WHERE login = $1", login)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) GetCharacter(id int64) (result model.Character, err error) {
	err = d.tx.Get(&result,
		`SELECT c.*, ac.account_id FROM characters c 
    JOIN account_characters ac on c.id = ac.character_id WHERE c.id = $1`, id)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) AddCharacter(name string) (id int, err error) {
	err = d.tx.Get(&id, "INSERT INTO characters VALUES (DEFAULT, $1, DEFAULT, DEFAULT) RETURNING id", name)
	return id, d.handleError(err)
}

func (d *DatabaseTransaction) DeleteCharacter(id int64) error {
	_, err := d.tx.Exec("DELETE FROM characters WHERE id = $1", id)
	return d.handleError(err)
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
