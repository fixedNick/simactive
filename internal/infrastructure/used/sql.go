package usedrepository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"
	"simactive/internal/lib/logger/sl"

	"github.com/go-sql-driver/mysql"
)

type UsedSQLRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUsedSQLRepository(db *sql.DB, logger *slog.Logger) *UsedSQLRepository {
	return &UsedSQLRepository{
		db:     db,
		logger: logger,
	}
}

func (ur *UsedSQLRepository) Add(ctx context.Context, simId int, serviceId int, isBlocked bool, blockedInfo string) (int, error) {
	const op = "UsedSQLRepository.Add"

	query := `INSERT INTO used_services (sim_id, service_id, is_blocked, blocked_info) VALUES (?, ?, ?, ?);`

	res, err := ur.db.ExecContext(ctx, query, simId, serviceId, isBlocked, blockedInfo)
	if err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ur.logger.Info(
				"Used service already exists",
				slog.String("op", op),
				slog.String("query", query),
				slog.Int("sim id", simId),
				slog.Int("service id", serviceId),
				slog.Bool("is blocked", isBlocked),
				slog.String("blocked info", blockedInfo),
			)
			return 0, repoerrors.ErrAlreadyExists
		}

		ur.logger.Error(
			"Failed to add used service",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("sim id", simId),
			slog.Int("service id", serviceId),
			slog.Bool("is blocked", isBlocked),
			slog.String("blocked info", blockedInfo),
			sl.Err(err),
		)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		ur.logger.Error(
			"Failed to get last insert id",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("sim id", simId),
			slog.Int("service id", serviceId),
			slog.Bool("is blocked", isBlocked),
			slog.String("blocked info", blockedInfo),
			sl.Err(err),
		)
		return 0, err
	}

	ur.logger.Info(
		"Used service successfully added",
		slog.String("op", op),
		slog.String("query", query),
		slog.Int("id", int(id)),
	)
	return int(id), nil
}
func (ur *UsedSQLRepository) GetList(ctx context.Context) (*core.List[*core.Used], error) {
	const op = "UsedSQLRepository.GetList"

	query := "SELECT id, sim_id, service_id, is_blocked, blocked_info FROM used_services"

	rows, err := ur.db.QueryContext(ctx, query)
	if err != nil {

		ur.logger.Error(
			"Failed to get used service list",
			slog.String("op", op),
			slog.String("query", query),
			sl.Err(err),
		)
		return nil, err
	}

	usedList := make(core.List[*core.Used], 0)
	for rows.Next() {
		var (
			id          int
			simId       int
			serviceId   int
			isBlocked   bool
			blockedInfo string
		)

		if err = rows.Scan(&id, &simId, &serviceId, &isBlocked, &blockedInfo); err != nil {
			ur.logger.Error(
				"Failed to scan used service",
				slog.String("op", op),
				slog.String("query", query),
				sl.Err(err),
			)
			return nil, err
		}
		used := core.NewUsed(id, simId, serviceId, isBlocked, blockedInfo)
		usedList[used.Id()] = &used
	}

	ur.logger.Info(
		"Used service list successfully got",
		slog.String("op", op),
		slog.String("query", query),
		slog.Int("count", len(usedList)),
	)
	return &usedList, nil

}
func (ur *UsedSQLRepository) ByID(ctx context.Context, id int) (*core.Used, error) {
	const op = "UsedSQLRepository.ByID"

	query := "SELECT sim_id, service_id, is_blocked, blocked_info FROM used_services WHERE id = ?"

	var (
		simId       int
		serviceId   int
		isBlocked   bool
		blockedInfo string
	)

	if err := ur.db.QueryRowContext(ctx, query, id).Scan(&simId, &serviceId, &isBlocked, &blockedInfo); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ur.logger.Info(
				"Used service does not exist",
				slog.String("op", op),
				slog.String("query", query),
				slog.Int("id", id),
			)
			return nil, repoerrors.ErrNotFound
		}
		ur.logger.Error(
			"Failed to get used service",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("id", id),
			sl.Err(err),
		)
		return nil, err
	}
	used := core.NewUsed(id, simId, serviceId, isBlocked, blockedInfo)
	ur.logger.Info(
		"Used service successfully got",
		slog.String("op", op),
		slog.String("query", query),
		slog.Int("id", id),
	)

	return &used, nil
}
func (ur *UsedSQLRepository) Update(ctx context.Context, s *core.Used) error {
	const op = "UsedSQLRepository.Update"

	query := "UPDATE used_services SET sim_id = ?, service_id = ?, is_blocked = ?, blocked_info = ? WHERE id = ?"

	_, err := ur.db.ExecContext(ctx, query, s.SimID(), s.ServiceID(), s.IsBlocked(), s.BlockedInfo(), s.Id())
	if err != nil {
		ur.logger.Error(
			"Failed to update used service",
			slog.String("op", op),
			slog.String("query", query),
			slog.Any("used", s),
			sl.Err(err),
		)
		return err
	}
	ur.logger.Info(
		"Used service successfully updated",
		slog.String("op", op),
		slog.String("query", query),
	)
	return nil
}
func (ur *UsedSQLRepository) Remove(ctx context.Context, id int) error {
	const op = "UsedSQLRepository.Remove"

	query := "DELETE FROM used_services WHERE id = ?"
	res, err := ur.db.ExecContext(ctx, query, id)
	if err != nil {
		ur.logger.Error(
			"Failed to remove used service",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("id", id),
			sl.Err(err),
		)
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		ur.logger.Error(
			"Failed to receive affected rows after query",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("id", id),
			sl.Err(err),
		)
		return err
	}

	if affectedRows == 0 {
		ur.logger.Info(
			"Used service does not exist",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("id", id),
		)
		return repoerrors.ErrNotFound
	}

	ur.logger.Info(
		"Used service successfully removed",
		slog.String("op", op),
		slog.String("query", query),
		slog.Int("id", id),
		slog.Int64("affected rows", affectedRows),
	)
	return nil
}
