package providerrepository

import (
	"context"
	"errors"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"
	"simactive/internal/lib/logger/sl"
)

type ProviderInMemory struct {
	logger *slog.Logger
	list   core.List[*core.Provider]
}

// NewProviderInMemory creates a new ProviderInMemory instance.
//
// It takes a logger as a parameter and returns a pointer to ProviderInMemory.
func NewProviderInMemory(logger *slog.Logger) *ProviderInMemory {
	return &ProviderInMemory{
		logger: logger,
		list:   make(core.List[*core.Provider]),
	}
}

// Add adds a new provider to the ProviderInMemory list.
//
// Parameters:
//
//	ctx context.Context - The context for the operation.
//	id int - The ID of the provider to add.
//	name string - The name of the provider to add.
//
// Return:
//
//	error - An error if the provider already exists, otherwise nil.
func (im *ProviderInMemory) Add(ctx context.Context, id int, name string) error {
	const op = "ProviderInMemory.Add"

	if provider, err := im.list.ByID(id); err == nil {

		im.logger.Info(
			"Provider already exists",
			slog.String("op", op),
			slog.Int("provider id", id),
			slog.String("provider name", name),
			slog.Any("provider", *provider),
		)
		return repoerrors.ErrAlreadyExists

	}

	p := core.NewProvider(id, name)
	im.list[id] = &p

	im.logger.Info(
		"Provider added",
		slog.String("op", op),
		slog.Int("provider id", id),
		slog.String("provider name", name),
	)
	return nil
}

// GetList retrieves the provider list.
//
// ctx context.Context
// *core.List[*core.Provider], error
func (im *ProviderInMemory) GetList(ctx context.Context) (*core.List[*core.Provider], error) {
	const op = "ProviderInMemory.GetList"

	im.logger.Info(
		"Provider list successfully retrieved",
		slog.String("op", op),
		slog.Int("provider count", len(im.list)),
	)
	return &im.list, nil
}

// ByID retrieves a provider by its ID.
//
// ctx context.Context, id int. Returns *core.Provider, error.
func (im *ProviderInMemory) ByID(ctx context.Context, id int) (*core.Provider, error) {
	const op = "ProviderInMemory.ByID"

	provider, err := im.list.ByID(id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			im.logger.Info(
				"Provider does not exist",
				slog.String("op", op),
				slog.Int("provider id", id),
			)
			return nil, repoerrors.ErrNotFound
		}

		im.logger.Error(
			"Failed to retrieve provider",
			slog.String("op", op),
			slog.Int("provider id", id),
			sl.Err(err),
		)
		return nil, err
	}

	im.logger.Info(
		"Provider successfully retrieved",
		slog.String("op", op),
		slog.Int("provider id", id),
		slog.Any("provider", *provider),
	)
	return provider, nil
}

// ByName retrieves a provider by name from the ProviderInMemory.
//
// ctx - the context
// name - the name of the provider
// *core.Provider - the retrieved provider
// error - error if the provider is not found
func (im *ProviderInMemory) ByName(ctx context.Context, name string) (*core.Provider, error) {
	const op = "ProviderInMemory.ByName"

	provider, exists := im.list.ContainsFunc(func(p *core.Provider) bool {
		return p.Name() == name
	})
	if !exists {

		im.logger.Info(
			"Provider does not exist",
			slog.String("op", op),
			slog.String("provider name", name),
		)
		return nil, repoerrors.ErrNotFound
	}

	im.logger.Info(
		"Provider successfully retrieved",
		slog.String("op", op),
		slog.String("provider name", name),
	)
	return provider, nil
}

// Remove removes a provider from the ProviderInMemory list by its ID.
// It takes a context and an integer ID as parameters and returns an error.
func (im *ProviderInMemory) Remove(ctx context.Context, id int) error {
	const op = "ProviderInMemory.Remove"

	_, err := im.list.ByID(id)

	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			im.logger.Info(
				"Provider does not exist",
				slog.String("op", op),
				slog.Int("provider id", id),
			)
			return err
		}

		im.logger.Error(
			"Failed to retrieve provider",
			slog.String("op", op),
			slog.Int("provider id", id),
			sl.Err(err),
		)
		return err
	}

	delete(im.list, id)

	im.logger.Info(
		"Provider successfully removed",
		slog.String("op", op),
		slog.Int("provider id", id),
	)
	return nil
}
