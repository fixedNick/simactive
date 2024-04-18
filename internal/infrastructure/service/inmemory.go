package servicerepository

import (
	"context"
	"errors"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"
	"simactive/internal/lib/logger/sl"
)

type ServiceInMemory struct {
	list   core.List[*core.Service]
	logger *slog.Logger
}

func NewServiceInMemoryRepository(logger *slog.Logger) *ServiceInMemory {
	return &ServiceInMemory{
		list:   make(core.List[*core.Service]),
		logger: logger,
	}
}

// Add adds a new service to the ServiceInMemory list if it doesn't already exist.
//
// ctx: the context.Context for the operation.
// serviceId: the ID of the service to be added.
// name: the name of the service to be added.
// error: returns an error if the service already exists.
func (si *ServiceInMemory) Add(ctx context.Context, serviceId int, name string) error {
	const op = "ServiceInMemory.Add"

	if service, err := si.list.ByID(serviceId); err == nil {
		si.logger.Info(
			"Service already exists",
			slog.String("op", op),
			slog.Int("service id", serviceId),
			slog.Any("service", *service),
		)
		return repoerrors.ErrAlreadyExists
	}

	s := core.NewService(serviceId, name)
	si.list[serviceId] = &s

	si.logger.Info(
		"Service added in memory",
		slog.String("op", op),
		slog.Int("service id", serviceId),
		slog.String("service name", name),
	)
	return nil
}

// Remove removes a service from memory by ID.
//
// ctx: context.Context
// id: int - the ID of the service to remove
// error - returns an error if the service with the given ID is not found
func (si *ServiceInMemory) Remove(ctx context.Context, id int) error {
	const op = "ServiceInMemory.Remove"

	service, err := si.list.ByID(id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			si.logger.Info(
				"Service does not exist",
				slog.String("op", op),
				slog.Int("service id", id),
			)
			return repoerrors.ErrNotFound
		}

		si.logger.Error(
			"Failed to retrieve service",
			slog.String("op", op),
			slog.Int("service id", id),
			sl.Err(err),
		)
		return err
	}

	delete(si.list, id)

	si.logger.Info(
		"Service removed from memory",
		slog.String("op", op),
		slog.Int("service id", id),
		slog.Any("removed service", *service),
	)
	return nil
}

// GetList retrieves the list of services from memory.
//
// ctx: context.Context
// Returns a pointer to a list of core.Service and an error.
func (si *ServiceInMemory) GetList(ctx context.Context) (*core.List[*core.Service], error) {
	const op = "ServiceInMemory.GetList"

	si.logger.Info(
		"Service list successfully retrieved",
		slog.String("op", op),
		slog.Int("service count", len(si.list)),
	)
	return &si.list, nil
}

// Update updates a service in memory.
//
// ctx: context.Context - The context for the operation.
// s: *core.Service - The service to update.
// error - Returns an error if the service is not found.
func (si *ServiceInMemory) Update(ctx context.Context, s *core.Service) error {
	const op = "ServiceInMemory.Update"

	_, err := si.list.ByID(s.Id())
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			si.logger.Info(
				"Service does not exist",
				slog.String("op", op),
				slog.Int("service id", s.Id()),
			)
			return repoerrors.ErrNotFound
		}

		si.logger.Error(
			"Failed to retrieve service",
			slog.String("op", op),
			slog.Int("service id", s.Id()),
			sl.Err(err),
		)
		return err
	}

	si.list[s.Id()] = s

	si.logger.Info(
		"Service successfully updated",
		slog.String("op", op),
		slog.Int("service id", s.Id()),
		slog.Any("updated values", *s),
	)
	return nil
}
