package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"simactive/internal/core"
)

type SimSqlRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

// Creates a new SQL Repository object
func NewSimSQLRepository(db *sql.DB, logger *slog.Logger) *SimSqlRepository {
	return &SimSqlRepository{
		db:     db,
		logger: logger,
	}
}

// Saving [s core.Sim] into database
// Returns [id] of saved core
// Return error:
// 1. Sim already exist
// 2. Internal error: Executing sql or fetching last insert id
func (sqlRepo *SimSqlRepository) Save(ctx context.Context, s core.Sim) (id int, err error) {
	const op = "sqlRepo.Save()"

	query := "INSERT INTO sim VALUES (?, ?, ?, ?, ?, ?)"
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

// Removes [s core.Sim] from database using [s.id]
// Return only internal (sql) errors
func (sqlRepo *SimSqlRepository) Remove(ctx context.Context, s core.Sim) error {
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

// Receiving list [core.SimList] from db
// Return only internal (sql) errors
func (sqlRepo *SimSqlRepository) GetSimList(ctx context.Context) (*core.SimList, error) {
	const op = "sqlRepo.SimList()"

	query := `SELECT * FROM sim`
	rows, err := sqlRepo.db.Query(query)

	if err != nil {
		panic("implement")
	}

	coreList := core.NewSimList()
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

		coreList[number] = core.NewSim(id, number, providerId, isActivated, activateUntil, isBlocked)
	}

	return &coreList, nil
}

// Gets [s Sim] from db by its own [id]
// Return error sql.ErrNoRows
// Return internal (sql) errors
func (sqlRepo *SimSqlRepository) ByID(ctx context.Context, id int) (core.Sim, error) {
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

	return core.NewSim(id, number, providerId, isActivated, activateUntil, isBlocked), nil
}

// Gets [s Sim] from db by its own [number]
// Return error sql.ErrNoRows
// Return internal (sql) errors
func (sqlRepo *SimSqlRepository) ByNumber(ctx context.Context, number string) (core.Sim, error) {
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

	return core.NewSim(id, number, providerId, isActivated, activateUntil, isBlocked), nil
}
