package usedrepository

import (
	"context"
	"errors"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"
	"simactive/internal/lib/logger/sl"
)

type UsedInMemoryRepository struct {
	list   core.List[*core.Used]
	logger *slog.Logger
}

func NewUsedInMemoryRepository(logger *slog.Logger) *UsedInMemoryRepository {
	return &UsedInMemoryRepository{
		list:   make(core.List[*core.Used]),
		logger: logger,
	}
}

func (ir *UsedInMemoryRepository) Add(ctx context.Context, id int, simId int, serviceId int, isBlocked bool, blockedInfo string) error {
	const op = "UsedInMemoryRepository.Add"

	if _, err := ir.list.ByID(id); err == nil {

		ir.logger.Info(
			"Used already exists",
			slog.String("op", op),
			slog.Int("id", id),
			slog.Int("sim id", simId),
			slog.Int("service id", serviceId),
			slog.Bool("is blocked", isBlocked),
			slog.String("blocked info", blockedInfo),
		)
		return repoerrors.ErrAlreadyExists
	}

	used := core.NewUsed(id, simId, serviceId, isBlocked, blockedInfo)
	ir.list[used.Id()] = &used

	ir.logger.Info(
		"Used added",
		slog.String("op", op),
		slog.Int("id", id),
		slog.Int("sim id", simId),
		slog.Int("service id", serviceId),
		slog.Bool("is blocked", isBlocked),
		slog.String("blocked info", blockedInfo),
	)
	return nil
}
func (ir *UsedInMemoryRepository) GetList(ctx context.Context) (*core.List[*core.Used], error) {
	const op = "UsedInMemoryRepository.GetList"

	ir.logger.Info(
		"Used list successfully retrieved",
		slog.String("op", op),
		slog.Int("used count", len(ir.list)),
	)
	return &ir.list, nil
}
func (ir *UsedInMemoryRepository) ByID(ctx context.Context, id int) (*core.Used, error) {
	const op = "UsedInMemoryRepository.ByID"

	used, err := ir.list.ByID(id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			ir.logger.Info(
				"Used does not exist",
				slog.String("op", op),
				slog.Int("used id", id),
			)
			return nil, err
		}

		ir.logger.Error(
			"Failed to retrieve used",
			slog.String("op", op),
			slog.Int("used id", id),
			sl.Err(err),
		)
		return nil, err
	}

	ir.logger.Info(
		"Used successfully retrieved",
		slog.String("op", op),
		slog.Int("used id", id),
	)
	return used, nil
}
func (ir *UsedInMemoryRepository) Update(ctx context.Context, s *core.Used) error {
	const op = "UsedInMemoryRepository.Update"

	if _, err := ir.list.ByID(s.Id()); err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			ir.logger.Info(
				"Used does not exist",
				slog.String("op", op),
				slog.Int("used id", s.Id()),
			)
			return err
		}

		ir.logger.Error(
			"Failed to retrieve used",
			slog.String("op", op),
			slog.Int("used id", s.Id()),
			sl.Err(err),
		)
		return err
	}

	ir.list[s.Id()] = s

	ir.logger.Info(
		"Used successfully updated",
		slog.String("op", op),
		slog.Int("used id", s.Id()),
	)

	return nil
}
func (ir *UsedInMemoryRepository) Remove(ctx context.Context, id int) error {
	const op = "UsedInMemoryRepository.Remove"

	_, err := ir.list.ByID(id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			ir.logger.Info(
				"Used does not exist",
				slog.String("op", op),
				slog.Int("used id", id),
			)
			return err
		}

		ir.logger.Error(
			"Failed to retrieve used",
			slog.String("op", op),
			slog.Int("used id", id),
			sl.Err(err),
		)
		return err
	}

	delete(ir.list, id)

	ir.logger.Info(
		"Used successfully removed",
		slog.String("op", op),
		slog.Int("used id", id),
	)

	return nil
}
