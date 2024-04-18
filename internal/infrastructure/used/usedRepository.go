package usedrepository

import (
	"context"
	"database/sql"
	"log/slog"
	"simactive/internal/core"
)

type UsedInMemory interface {
	SamemRepoFuncs
	Add(ctx context.Context, id int, simId int, serviceId int, isBlocked bool, blockedInfo string) error
}

type UsedSQL interface {
	SamemRepoFuncs
	Add(ctx context.Context, simId int, serviceId int, isBlocked bool, blockedInfo string) (id int, err error)
}

type SamemRepoFuncs interface {
	GetList(ctx context.Context) (*core.List[*core.Used], error)
	ByID(ctx context.Context, id int) (*core.Used, error)
	Update(ctx context.Context, s *core.Used) error
	Remove(ctx context.Context, id int) error
}

type UsedRepository struct {
	logger   *slog.Logger
	db       *sql.DB
	inMemory UsedInMemory
	sql      UsedSQL
}

func NewUsedRepository(logger *slog.Logger, db *sql.DB, inMemory UsedInMemory, sql UsedSQL) *UsedRepository {
	const op = "repository.used.NewUsedRepository"

	logger.Info("Used Repository initialized", slog.String("op", op))

	return &UsedRepository{
		logger:   logger,
		db:       db,
		inMemory: inMemory,
		sql:      sql,
	}
}

func (ur *UsedRepository) Add(ctx context.Context, simId int, serviceId int, isBlocked bool, blockedInfo string) (int, error) {

	id, err := ur.sql.Add(ctx, simId, serviceId, isBlocked, blockedInfo)

	if err != nil {
		return 0, err
	}

	if err = ur.inMemory.Add(ctx, id, simId, serviceId, isBlocked, blockedInfo); err != nil {
		return 0, err
	}
	return id, nil
}
func (ur *UsedRepository) GetList(ctx context.Context) (*core.List[*core.Used], error) {
	list, err := ur.inMemory.GetList(ctx)
	if err != nil {
		return nil, err
	}

	if list != nil && len(*list) > 0 {
		return list, nil
	}

	list, err = ur.sql.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}
func (ur *UsedRepository) ByID(ctx context.Context, id int) (*core.Used, error) {
	if used, err := ur.inMemory.ByID(ctx, id); err == nil {
		return used, nil
	}

	return ur.sql.ByID(ctx, id)
}
func (ur *UsedRepository) Update(ctx context.Context, s *core.Used) error {
	if err := ur.inMemory.Update(ctx, s); err != nil {
		return err
	}

	if err := ur.sql.Update(ctx, s); err != nil {
		return err
	}

	return nil
}
func (ur *UsedRepository) Remove(ctx context.Context, id int) error {
	err := ur.inMemory.Remove(ctx, id)
	if err != nil {
		return err
	}

	err = ur.sql.Remove(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
