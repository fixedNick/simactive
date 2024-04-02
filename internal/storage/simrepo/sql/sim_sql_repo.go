package simrepo

import (
	"database/sql"
	"log/slog"
	"simactive/internal/core/sim"
)

type SqlRepository interface {
	Save(s sim.Sim) (id int, err error)
	Remove(s sim.Sim) error
	SimList() (*sim.SimList, error)
	ByID(id int) (sim.Sim, error)
	ByNumber(number string) (sim.Sim, error)
}

type SimSqlRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// Creates a new SQL Repository object
func NewSQLRepository(db *sql.DB, logger *slog.Logger) *SimSqlRepository {
	return &SimSqlRepository{
		db:     db,
		logger: logger,
	}
}

// Saving [s sim.Sim] into database
// Returns [id] of saved sim
// Return error:
// 1. Sim already exist
// 2. Internal error: Executing sql or fetching last insert id
func (sqlRepo *SimSqlRepository) Save(s sim.Sim) (id int, err error) {
	const op = "sqlRepo.Save()"

	query := "INSERT INTO sim (?, ?, ?, ?, ?, ?)"
	r, err := sqlRepo.db.Exec(query, 0, s.Number(), s.ProviderID(), s.IsActivated(), s.ActivateUntil(), s.IsBlocked())
	if err != nil {
		panic("implement")
	}

	insertedId, err := r.LastInsertId()
	if err != nil {
		panic("implement")
	}

	return int(insertedId), nil
}

// Removes [s sim.Sim] from database using [s.id]
// Return only internal (sql) errors
func (sqlRepo *SimSqlRepository) Remove(s sim.Sim) error {
	const op = "sqlRepo.Remove()"

	query := "DELETE FROM sim WHERE id = ?"
	r, err := sqlRepo.db.Exec(query, s.Id())
	if err != nil {
		panic("implement")
	}

	affected, err := r.RowsAffected()
	if err != nil {
		// заглушка на affected
		panic("implement" + string(affected))
	}
	return nil
}

// Receiving list [sim.SimList] from db
// Return only internal (sql) errors
func (sqlRepo *SimSqlRepository) SimList() (*sim.SimList, error) {
	const op = "sqlRepo.SimList()"

	query := `SELECT * FROM sim`
	rows, err := sqlRepo.db.Query(query)

	if err != nil {
		panic("implement")
	}

	simList := sim.NewSimList()
	for {
		if !rows.Next() {
			break
		}
		var (
			id            int
			number        string
			providerId    int
			isActivated   bool
			activateUntil int64
			isBlocked     bool
		)

		err = rows.Scan(&id, &number, &providerId, &isActivated, &activateUntil, &isBlocked)
		if err != nil {
			panic("implement")
		}

		simList[number] = sim.NewSim(id, number, providerId, isActivated, activateUntil, isBlocked)
	}

	return &simList, nil
}

// Gets [s Sim] from db by its own [id]
// Return error sql.ErrNoRows
// Return internal (sql) errors
func (sqlRepo *SimSqlRepository) ByID(id int) (sim.Sim, error) {
	const op = "sqlRepo.ById()"

	query := "SELECT number, provider_id, is_activated, activate_until, is_blocked FROM sim WHERE id = ?"
	row := sqlRepo.db.QueryRow(query, id)

	var (
		providerId             int
		number                 string
		isActivated, isBlocked bool
		activateUntil          int64
	)

	err := row.Scan(&number, &providerId, &isActivated, &activateUntil, &isBlocked)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		panic("implement no rows")
	default:
		panic("implement internal")
	}

	return sim.NewSim(id, number, providerId, isActivated, activateUntil, isBlocked), nil
}

// Gets [s Sim] from db by its own [number]
// Return error sql.ErrNoRows
// Return internal (sql) errors
func (sqlRepo *SimSqlRepository) ByNumber(number string) (sim.Sim, error) {
	const op = "sqlRepo.ByNumber()"

	query := "SELECT id, provider_id, is_activated, activate_until, is_blocked FROM sim WHERE number = ?"
	row := sqlRepo.db.QueryRow(query, number)

	var (
		providerId, id         int
		isActivated, isBlocked bool
		activateUntil          int64
	)

	err := row.Scan(&id, &providerId, &isActivated, &activateUntil, &isBlocked)

	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		panic("implement no rows")
	default:
		panic("implement internal")
	}

	return sim.NewSim(id, number, providerId, isActivated, activateUntil, isBlocked), nil
}
