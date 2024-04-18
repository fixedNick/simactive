package simrepository

import (
	"context"
	"errors"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"
	"simactive/internal/lib/logger/sl"
)

// SimInMemory is a repository that stores SIM cards in memory.
type SimInMemory struct {
	list   core.List[*core.Sim]
	logger *slog.Logger
}

func NewSimInMemoryRepository(logger *slog.Logger) *SimInMemory {
	return &SimInMemory{
		list:   make(core.List[*core.Sim]),
		logger: logger,
	}
}

// Add adds a new SIM card to the SimInMemory instance.
//
// Parameters:
//   - simId: the ID of the SIM card
//   - number: the phone number associated with the SIM card
//   - provider: the provider of the SIM card
//   - isActivated: flag indicating if the SIM card is activated
//   - activateUntil: timestamp until the SIM card is activated
//   - isBlocked: flag indicating if the SIM card is blocked
//
// Return:
//   - err: an error, if any
//   - ErrAlreadyExists: if the SIM card already exists
func (i *SimInMemory) Add(ctx context.Context, simId int, number string, provider *core.Provider, isActivated bool, activateUntil int64, isBlocked bool) (err error) {
	const op = "SimInMemory.Add"

	if sim, err := i.list.ByID(simId); err == nil {

		i.logger.Info(
			"Sim already exists",
			slog.String("op", op),
			slog.Any("sim", *sim),
			slog.Int("sim id", simId),
			slog.String("number", number),
			slog.Int("provider id", provider.Id()),
			slog.String("provider name", provider.Name()),
			slog.Bool("isActivated", isActivated),
			slog.Int64("activateUntil", activateUntil),
			slog.Bool("isBlocked", isBlocked),
		)

		return repoerrors.ErrAlreadyExists
	}

	s := core.NewSim(simId, number, provider, isActivated, activateUntil, isBlocked)
	i.list[simId] = &s

	i.logger.Info(
		"Sim successfully added",
		slog.String("op", op),
		slog.Int("sim id", simId),
		slog.String("number", number),
		slog.Int("provider id", provider.Id()),
		slog.String("provider name", provider.Name()),
		slog.Bool("isActivated", isActivated),
		slog.Int64("activateUntil", activateUntil),
		slog.Bool("isBlocked", isBlocked),
	)
	return nil
}

// Remove removes a Sim from the SimInMemory list.
//
// ctx context.Context, id int
// error
func (i *SimInMemory) Remove(ctx context.Context, id int) error {
	const op = "SimInMemory.Remove"

	_, err := i.list.ByID(id)

	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			i.logger.Info(
				"Sim does not exist",
				slog.String("op", op),
				slog.Int("sim id", id),
			)
			return err
		}

		i.logger.Error("Failed to retrieve sim",
			slog.String("op", op),
			slog.Int("sim id", id),
			sl.Err(err),
		)
		return err
	}

	delete(i.list, id)

	i.logger.Info(
		"Sim successfully removed",
		slog.String("op", op),
		slog.Int("sim id", id),
	)
	return nil
}

// GetList retrieves the list of Sims.
//
// Context ctx - The context for the operation.
// Returns a pointer to List of Sims and an error.
func (i *SimInMemory) GetList(ctx context.Context) (*core.List[*core.Sim], error) {
	const op = "SimInMemory.GetList"

	i.logger.Info(
		"Sim list successfully retrieved",
		slog.String("op", op),
		slog.Int("sim count", len(i.list)),
	)
	return &i.list, nil
}

// Update updates the Sim in the SimInMemory with the given context and core.Sim.
//
// ctx context.Context, s *core.Sim
// error
func (i *SimInMemory) Update(ctx context.Context, s *core.Sim) error {
	const op = "SimInMemory.Update"

	_, err := i.list.ByID(s.Id())
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			i.logger.Info(
				"Sim does not exist",
				slog.String("op", op),
				slog.Int("sim id", s.Id()),
			)
			return err
		}

		i.logger.Error("Failed to retrieve sim",
			slog.String("op", op),
			slog.Int("sim id", s.Id()),
			sl.Err(err),
		)
		return err
	}

	i.list[s.Id()] = s

	i.logger.Info(
		"Sim successfully updated",
		slog.String("op", op),
		slog.Int("sim id", s.Id()),
		slog.Any("updated values", *s),
	)
	return nil
}

func (i *SimInMemory) ByID(ctx context.Context, id int) (*core.Sim, error) {
	const op = "SimInMemory.ByID"

	sim, err := i.list.ByID(id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			i.logger.Info(
				"Sim does not exist",
				slog.String("op", op),
				slog.Int("sim id", id),
			)
			return nil, err
		}

		i.logger.Error("Failed to retrieve sim",
			slog.String("op", op),
			slog.Int("sim id", id),
			sl.Err(err),
		)
		return nil, err
	}

	i.logger.Info(
		"Sim successfully retrieved",
		slog.String("op", op),
		slog.Int("sim id", id),
		slog.Any("sim", *sim),
	)
	return sim, nil
}
