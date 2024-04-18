package simrepository

import (
	"context"
	"database/sql"
	"log/slog"
	"simactive/internal/core"
)

type SimInMemRepo interface {
	SameRepoFuncs
	Add(ctx context.Context, simId int, number string, provider *core.Provider, isActivated bool, activateUntil int64, isBlocked bool) (err error)
}

type SimSQLRepo interface {
	SameRepoFuncs
	Add(ctx context.Context, number string, provider *core.Provider, isActivated bool, activateUntil int64, isBlocked bool) (simId int, err error)
}

type SameRepoFuncs interface {
	Remove(ctx context.Context, id int) (err error)
	GetList(ctx context.Context) (*core.List[*core.Sim], error)
	Update(ctx context.Context, s *core.Sim) error
	ByID(ctx context.Context, id int) (*core.Sim, error)
}

type SimRepository struct {
	logger   *slog.Logger
	db       *sql.DB
	inMemory SimInMemRepo
	sql      SimSQLRepo
}

// NewRepository initializes a new Repository with the given logger, database, in-memory repository, and SQL repository.
//
// Parameters:
// - logger: a slog.Logger instance for logging
// - db: a *sql.DB instance for database operations
// - simInMemory: a SimInMemRepo instance for in-memory repository operations
// - simSQL: a SimSQLRepo instance for SQL repository operations
// Return type: *Repository
func NewSimRepository(logger *slog.Logger, db *sql.DB, simInMemory SimInMemRepo, simSQL SimSQLRepo) *SimRepository {
	const op = "repository.sim.NewRepository"

	logger.Info("Sim Repository initialized", slog.String("op", op))
	return &SimRepository{
		logger:   logger,
		db:       db,
		inMemory: simInMemory,
		sql:      simSQL,
	}
}

// Add adds a new sim into in-memory and into sql
// If errors not occured it will return [ID] of new sim
func (r *SimRepository) Add(ctx context.Context, number string, provider *core.Provider, isActivated bool, activateUntil int64, isBlocked bool) (int, error) {
	id, err := r.sql.Add(ctx, number, provider, isActivated, activateUntil, isBlocked)
	if err != nil {
		return 0, err
	}

	err = r.inMemory.Add(ctx, id, number, provider, isActivated, activateUntil, isBlocked)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Remove removes a sim with the given id.
//
// ctx: the context in which the operation is performed.
// id: the id of the simulation to be removed.
// error: an error if any occurred during the removal process.
// Possibly errors is repository.ErrNotFound if sim with given id does not exist.
func (r *SimRepository) Remove(ctx context.Context, id int) (err error) {
	err = r.inMemory.Remove(ctx, id)
	if err != nil {
		return err
	}

	err = r.sql.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

// GetList retrieves a list of sims from the repository.
//
// ctx context.Context
// *core.List[*core.Sim], error
func (r *SimRepository) GetList(ctx context.Context) (*core.List[*core.Sim], error) {
	list, err := r.inMemory.GetList(ctx)
	if err != nil {
		return nil, err
	}

	if list != nil && len(*list) != 0 {
		return list, nil
	}

	list, err = r.sql.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Update updates the SimRepository with the given Sim.
//
// ctx context.Context, s *core.Sim
// error
func (r *SimRepository) Update(ctx context.Context, s *core.Sim) error {
	if err := r.inMemory.Update(ctx, s); err != nil {
		return err
	}

	if err := r.sql.Update(ctx, s); err != nil {
		return err
	}

	return nil
}

func (r *SimRepository) ByID(ctx context.Context, id int) (*core.Sim, error) {
	if s, err := r.inMemory.ByID(ctx, id); err == nil {
		return s, nil
	}

	return r.sql.ByID(ctx, id)
}
