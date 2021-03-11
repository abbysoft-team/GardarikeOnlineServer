package postgres

import (
	"abbysoft/gardarike-online/db"
	"abbysoft/gardarike-online/model"
	rpc "abbysoft/gardarike-online/rpc/generated"
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

type allBuildingsRow struct {
	CharacterID int64  `db:"character_id"`
	BuildingID  int64  `db:"building_id"`
	Count       uint64 `db:"count"`
}

func (d *DatabaseTransaction) UpdateProductionRates(rates model.Resources) error {
	_, err := d.tx.NamedExec(`UPDATE production_rates SET wood=:wood, leather=:leather, food=:food, stone=:stone
WHERE character_id=:character_id`, rates)
	return d.handleError(err)
}

func (d *DatabaseTransaction) GetProductionRates(characterID int64) (result model.Resources, err error) {
	err = d.tx.Get(&result, "SELECT * FROM production_rates WHERE character_id=$1", characterID)
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) GetAllBuildings() (result map[int64]model.CharacterBuildings, err error) {
	var rows []allBuildingsRow
	err = d.tx.Select(&rows, `select c.id character_id, tb.building_id, COUNT(tb.building_id) from town_buildings tb 
join towns t on tb.town_id = t.id 
join characters c on t.owner_name = c.name
GROUP BY c.id, tb.building_id`)

	if err != nil {
		return nil, d.handleError(err)
	}

	result = make(map[int64]model.CharacterBuildings)
	for _, row := range rows {
		if result[row.CharacterID] == nil {
			result[row.CharacterID] = make(model.CharacterBuildings)
		}

		if model.IsValidBuildingType(int32(row.BuildingID)) {
			result[row.CharacterID][rpc.BuildingType(row.BuildingID)] = row.Count
		}
	}

	return
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

func (d *DatabaseTransaction) AddTownBuilding(townID int64, building model.Building) error {
	_, err := d.tx.Exec("INSERT INTO town_buildings VALUES ($1, $2, $3, $4)",
		townID, building.ID, int64(building.Location.X), int64(building.Location.Y))
	return d.handleError(err)
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
		`INSERT INTO towns VALUES (DEFAULT, :x, :y, :owner_name, :population, :name)`, town)
	return d.handleError(err)
}

func (d *DatabaseTransaction) UpdateResources(resources model.Resources) error {
	_, err := d.tx.NamedExec(
		`UPDATE resources SET stone = :stone, food = :food, leather = :leather, wood = :wood
WHERE character_id=:character_id`, resources)
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
	if err != nil {
		return d.handleError(err)
	}

	if err := d.UpdateResources(character.Resources); err != nil {
		return d.handleError(err)
	}

	err = d.UpdateProductionRates(character.ProductionRate)
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

	if err != nil {
		return result, d.handleError(err)
	}

	resources, err := d.GetResources(id)
	if err != nil {
		return result, fmt.Errorf("failed to get character resources: %w", err)
	}

	result.Resources = resources

	productionRates, err := d.GetProductionRates(id)
	if err != nil {
		return result, fmt.Errorf("failed to get character production rates: %w", err)
	}

	result.ProductionRate = productionRates
	return result, d.handleError(err)
}

func (d *DatabaseTransaction) AddCharacter(name string) (id int, err error) {
	err = d.tx.Get(&id, "INSERT INTO characters VALUES (DEFAULT, $1, DEFAULT, DEFAULT) RETURNING id", name)
	if err != nil {
		return 0, d.handleError(err)
	}

	_, err = d.tx.Exec("INSERT INTO resources VALUES ($1, 0, 0, 0, 0)", id)
	if err != nil {
		return id, fmt.Errorf("failed to insert character resources: %w", err)
	}

	_, err = d.tx.Exec("INSERT INTO production_rates VALUES ($1, 0, 0, 0, 0)", id)
	return id, fmt.Errorf("failed to insert character production rates: %w", err)
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
