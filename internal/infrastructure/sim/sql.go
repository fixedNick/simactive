package simrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"

	"github.com/go-sql-driver/mysql"
)

type SimSQL struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewSimSQLRepository(db *sql.DB, logger *slog.Logger) *SimSQL {
	return &SimSQL{
		db:     db,
		logger: logger,
	}
}

// Add adds a new Sim entry to the database.
//
// Parameters:
//   - ctx: the context of the operation
//   - number: the number associated with the Sim
//   - provider: the provider of the Sim
//   - isActivated: a boolean indicating if the Sim is activated
//   - activateUntil: the activation time of the Sim
//   - isBlocked: a boolean indicating if the Sim is blocked
//
// Returns:
//   - int: the ID of the inserted Sim
//   - error: an error, if any
func (ss *SimSQL) Add(ctx context.Context, number string, provider *core.Provider, isActivated bool, activateUntil int64, isBlocked bool) (int, error) {
	const op = "SimSQL.Add"

	query := `CALL InsertSim(?,?,?,?,?,@lastId);`

	var insertedId int

	_, err := ss.db.ExecContext(ctx, query, number, provider.Id(), isActivated, activateUntil, isBlocked)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ss.logger.Info(
				"Sim already exists",
				slog.String("op", op),
				slog.String("number", number),
				slog.Int("provider id", provider.Id()),
				slog.Bool("isActivated", isActivated),
				slog.Int64("activateUntil", activateUntil),
				slog.Bool("isBlocked", isBlocked),
			)

			return 0, repoerrors.ErrAlreadyExists
		}

		ss.logger.Info(
			"Failed to add sim",
			slog.String("op", op),
			slog.String("query", query),
			"err", err,
		)
		return 0, err
	}
	if err := ss.db.QueryRowContext(ctx, "SELECT @lastId").Scan(&insertedId); err != nil {
		fmt.Println("inner error: ", err)
		// validate is that already exist erro from sql
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ss.logger.Info(
				"Sim already exists",
				slog.String("op", op),
				slog.String("number", number),
				slog.Int("provider id", provider.Id()),
				slog.Bool("isActivated", isActivated),
				slog.Int64("activateUntil", activateUntil),
				slog.Bool("isBlocked", isBlocked),
			)

			return 0, repoerrors.ErrAlreadyExists
		}

		return 0, err
	}

	ss.logger.Info(
		"Sim added",
		slog.String("op", op),
		slog.Int("sim id", insertedId),
		slog.String("number", number),
		slog.Int("provider id", provider.Id()),
		slog.Bool("isActivated", isActivated),
		slog.Int64("activateUntil", activateUntil),
		slog.Bool("isBlocked", isBlocked),
	)
	return insertedId, nil
}

// Remove removes a sim from the database.
//
// ctx: context for the operation.
// id: the ID of the sim to be removed.
// error: returns any error that occurred during the operation.
func (ss *SimSQL) Remove(ctx context.Context, id int) (err error) {
	const op = "SimSQL.Remove"

	query := "DELETE FROM sim WHERE id = ?"
	res, err := ss.db.ExecContext(ctx, query, id)
	if err != nil {
		ss.logger.Warn(
			"Failed to remove sim",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("sim id", id),
			"err", err,
		)
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		ss.logger.Warn(
			"Failed to receive affected rows after query",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("sim id", id),
			"err", err,
		)
		return err
	}
	if affectedRows == 0 {
		ss.logger.Info(
			"Sim does not exist",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("sim id", id),
		)
		return repoerrors.ErrNotFound
	}

	ss.logger.Info(
		"Sim successfully removed",
		slog.String("op", op),
		slog.Int("sim id", id),
	)
	return nil
}

// GetList retrieves a list of Sims from the database.
//
// ctx: the context for the database query.
// Returns a list of Sims and an error if any.
func (ss *SimSQL) GetList(ctx context.Context) (*core.List[*core.Sim], error) {
	const op = "SimSQL.GetList"

	query := `SELECT sim.id, sim.number, sim.provider_id, sim.is_activated, sim.activate_until, sim.is_blocked, provider.name 
				FROM sim 
				JOIN provider
				ON provider.id = sim.provider_id`
	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		ss.logger.Warn(
			"Failed to get sim list",
			slog.String("op", op),
			slog.String("query", query),
			"err", err,
		)
	}

	simList := make(core.List[*core.Sim], 0)
	for rows.Next() {
		var (
			id            int
			number        string
			providerId    int
			isActivated   bool
			activateUntil int64
			isBlocked     bool
			providerName  string
		)

		err = rows.Scan(&id, &number, &providerId, &isActivated, &activateUntil, &isBlocked, &providerName)
		if err != nil {
			ss.logger.Warn(
				"Failed to scan sim",
				slog.String("op", op),
				slog.String("query", query),
				"err", err,
			)
			return nil, err
		}

		p := core.NewProvider(providerId, providerName)

		sim := core.NewSim(
			id,
			number,
			&p,
			isActivated,
			activateUntil,
			isBlocked,
		)
		simList[id] = &sim
	}

	ss.logger.Info(
		"Sim list successfully retrieved",
		slog.String("op", op),
		slog.Int("sim count", len(simList)),
	)
	return &simList, nil
}

// Update updates the Sim object in the database.
//
// ctx context.Context, s *core.Sim
// error
func (ss *SimSQL) Update(ctx context.Context, s *core.Sim) error {
	const op = "SimSQL.Update"

	logSimInfo := []any{
		slog.Int("sim id", s.Id()),
		slog.Int("provider id", s.Provider().Id()),
		slog.String("provider name", s.Provider().Name()),
		slog.String("number", s.Number()),
		slog.Bool("isActivated", s.IsActivated()),
		slog.Int64("activateUntil", s.ActivateUntil()),
		slog.Bool("isBlocked", s.IsBlocked()),
	}

	query := "UPDATE sim SET number = ?, provider_id = ?, is_activated = ?, activate_until = ?, is_blocked = ? WHERE id = ?"
	_, err := ss.db.ExecContext(ctx, query, s.Number(), s.Provider().Id(), s.IsActivated(), s.ActivateUntil(), s.IsBlocked(), s.Id())
	if err != nil {

		// TODO:
		// possibly unique constraint error

		ss.logger.Warn(
			"Failed to update sim",
			slog.String("op", op),
			slog.String("query", query),
			slog.Group("sim info", logSimInfo...),
			"err", err,
		)
		return err
	}

	ss.logger.Info(
		"Sim successfully updated",
		slog.String("op", op),
		slog.Group("sim info", logSimInfo...),
	)
	return nil
}

// ByID retrieves a Sim by its ID.
//
// Takes in a context and an integer ID, returns a pointer to a core.Sim and an error.
func (ss *SimSQL) ByID(ctx context.Context, id int) (*core.Sim, error) {
	const op = "SimSQL.ByID"

	query := `SELECT sim.id, sim.number, sim.provider_id, sim.is_activated, sim.activate_until, sim.is_blocked, provider.name 
				FROM sim 
				JOIN provider
				ON provider.id = sim.provider_id
				WHERE sim.id = ?`

	var (
		number        string
		providerId    int
		isActivated   bool
		activateUntil int64
		isBlocked     bool
		providerName  string
	)
	err := ss.db.QueryRowContext(ctx, query, id).Scan(&id, &number, &providerId, &isActivated, &activateUntil, &isBlocked, &providerName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ss.logger.Info(
				"Sim does not exist",
				slog.String("op", op),
				slog.Int("sim id", id),
			)
			return nil, repoerrors.ErrNotFound
		}

		ss.logger.Warn(
			"Failed to get sim",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("sim id", id),
			"err", err,
		)
		return nil, err
	}

	p := core.NewProvider(providerId, providerName)

	sim := core.NewSim(
		id,
		number,
		&p,
		isActivated,
		activateUntil,
		isBlocked,
	)
	ss.logger.Info(
		"Sim successfully retrieved",
		slog.String("op", op),
		slog.Any("sim", sim),
	)
	return &sim, nil
}
