package providerrepository

import (
	"context"
	"database/sql"
	"log/slog"
	"simactive/internal/core"
)

type ProviderInMemoryRepo interface {
	SamemRepoFuncs
	Add(ctx context.Context, id int, name string) error
}

type ProviderSQLRepo interface {
	SamemRepoFuncs
	Add(ctx context.Context, name string) (int, error)
}

type SamemRepoFuncs interface {
	GetList(ctx context.Context) (*core.List[*core.Provider], error)
	ByID(ctx context.Context, id int) (*core.Provider, error)
	ByName(ctx context.Context, name string) (*core.Provider, error)
	Remove(ctx context.Context, id int) error
}

type ProviderRepository struct {
	logger   *slog.Logger
	db       *sql.DB
	inMemory ProviderInMemoryRepo
	sql      ProviderSQLRepo
}

// NewProviderRepository initializes a new ProviderRepository.
//
// It takes a logger, a database connection, an in-memory provider repository, and a SQL provider repository as parameters.
// It returns a pointer to ProviderRepository.
func NewProviderRepository(logger *slog.Logger, db *sql.DB, inMemory ProviderInMemoryRepo, sql ProviderSQLRepo) *ProviderRepository {
	const op = "repository.provider.NewProviderRepository"

	logger.Info("Provider Repository initialized", slog.String("op", op))

	return &ProviderRepository{
		logger:   logger,
		db:       db,
		inMemory: inMemory,
		sql:      sql,
	}
}

func (r *ProviderRepository) Add(ctx context.Context, name string) (int, error) {
	id, err := r.sql.Add(ctx, name)
	if err != nil {
		return 0, err
	}

	err = r.inMemory.Add(ctx, id, name)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetList retrieves a list of providers.
//
// Context parameter.
// Returns a list of providers and an error.
func (r *ProviderRepository) GetList(ctx context.Context) (*core.List[*core.Provider], error) {

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

// ByID retrieves a provider by its ID.
//
// ctx - the context in which the operation should be performed.
// id - the ID of the provider to retrieve.
// Returns a core.Provider and an error.
func (r *ProviderRepository) ByID(ctx context.Context, id int) (*core.Provider, error) {
	if p, err := r.inMemory.ByID(ctx, id); err == nil {
		return p, nil
	}
	return r.sql.ByID(ctx, id)
}

// ByName retrieves a provider by name.
//
// ctx context.Context, name string. Returns core.Provider, error.
func (r *ProviderRepository) ByName(ctx context.Context, name string) (*core.Provider, error) {
	if p, err := r.inMemory.ByName(ctx, name); err == nil {
		return p, nil
	}
	return r.sql.ByName(ctx, name)
}

// Remove removes an item using the given id.
//
// Takes a context.Context and an int as parameters.
// Returns an error.
func (r *ProviderRepository) Remove(ctx context.Context, id int) error {
	if err := r.inMemory.Remove(ctx, id); err != nil {
		return err
	}
	return r.sql.Remove(ctx, id)
}
