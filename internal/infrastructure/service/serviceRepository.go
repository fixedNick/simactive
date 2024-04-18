package servicerepository

import (
	"context"
	"database/sql"
	"log/slog"
	"simactive/internal/core"
)

type ServiceInMemRepo interface {
	SameRepoFuncs
	Add(ctx context.Context, serviceId int, name string) (err error)
}

type ServiceSQLRepo interface {
	SameRepoFuncs
	Add(ctx context.Context, name string) (serviceId int, err error)
}

type SameRepoFuncs interface {
	Remove(ctx context.Context, id int) (err error)
	GetList(ctx context.Context) (*core.List[*core.Service], error)
	Update(ctx context.Context, s *core.Service) error
}

type ServiceRepository struct {
	logger   *slog.Logger
	db       *sql.DB
	inMemory ServiceInMemRepo
	sql      ServiceSQLRepo
}

func NewServiceRepository(logger *slog.Logger, db *sql.DB, serviceInMemory ServiceInMemRepo, serviceSQL ServiceSQLRepo) *ServiceRepository {
	const op = "repository.service.NewServiceRepository"
	logger.Info("Service Repository initialized", slog.String("op", op))
	return &ServiceRepository{
		logger:   logger,
		db:       db,
		inMemory: serviceInMemory,
		sql:      serviceSQL,
	}
}

// Add adds a new service to the repository.
//
// ctx: the context for the operation.
// name: the name of the service to add.
// Returns the service ID and an error if any.  Possibly errors: repository.ErrAlreadyExists.
func (sr *ServiceRepository) Add(ctx context.Context, name string) (serviceId int, err error) {

	id, err := sr.sql.Add(ctx, name)
	if err != nil {
		return 0, err
	}

	err = sr.inMemory.Add(ctx, id, name)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Remove removes an item with the given ID from the ServiceRepository.
//
// ctx: context.Context - The context for the operation.
// id: int - The ID of the item to be removed.
// error - Returns an error if any occurred during the removal process.  Possibly errors: repository.ErrNotFound.
func (sr *ServiceRepository) Remove(ctx context.Context, id int) (err error) {

	if err = sr.inMemory.Remove(ctx, id); err != nil {
		return err
	}

	if err = sr.sql.Remove(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetList retrieves a list of services from the ServiceRepository.
//
// ctx - the context for the operation.
// Returns a list of services and an error, if any.
func (sr *ServiceRepository) GetList(ctx context.Context) (*core.List[*core.Service], error) {
	list, err := sr.inMemory.GetList(ctx)
	if err != nil {
		return nil, err
	}

	if list != nil && len(*list) != 0 {
		return list, nil
	}

	list, err = sr.sql.GetList(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Update updates the service in the ServiceRepository.
//
// ctx: the context.Context for the operation.
// s: the *core.Service to be updated.
// error: returns an error if the update operation fails. Possibly errors: repository.ErrNotFound
func (sr *ServiceRepository) Update(ctx context.Context, s *core.Service) error {
	if err := sr.inMemory.Update(ctx, s); err != nil {
		return err
	}

	if err := sr.sql.Update(ctx, s); err != nil {
		return err
	}

	return nil
}
