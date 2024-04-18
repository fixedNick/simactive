package servicerepository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"simactive/internal/core"
	"simactive/internal/infrastructure/repoerrors"

	"github.com/go-sql-driver/mysql"
)

type ServiceSQL struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewServiceSQLRepository(db *sql.DB, logger *slog.Logger) *ServiceSQL {
	return &ServiceSQL{
		logger: logger,
		db:     db,
	}
}

// Add adds a new service with the given name to the database.
//
// ctx: the context for the request
// name: the name of the service to be added
// Returns the ID of the newly added service and an error, if any
func (ss *ServiceSQL) Add(ctx context.Context, name string) (int, error) {
	const op = "ServiceSQL.Add"

	query := "INSERT INTO service (name) VALUES (?)"
	res, err := ss.db.ExecContext(ctx, query, name)
	if err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {

			ss.logger.Info(
				"Service already exists",
				slog.String("op", op),
				slog.String("query", query),
				slog.String("service name", name),
			)

			return 0, repoerrors.ErrAlreadyExists
		}

		ss.logger.Warn(
			"Failed to add service",
			slog.String("op", op),
			slog.String("query", query),
			slog.String("service name", name),
			"err", err,
		)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {

		ss.logger.Warn(
			"Failed to receive last insert id after query",
			slog.String("op", op),
			slog.String("query", query),
			slog.String("service name", name),
			"err", err,
		)
		return 0, err
	}

	ss.logger.Info(
		"Service successfully added",
		slog.String("op", op),
		slog.String("service name", name),
		slog.Int64("service id", id),
	)
	return int(id), nil
}

// Remove removes a service from the database based on the provided ID.
// ctx - the context for the operation.
// id - the ID of the service to remove.
// error - returns an error if the removal operation encounters any issues.
func (ss *ServiceSQL) Remove(ctx context.Context, id int) error {
	const op = "ServiceSQL.Remove"

	query := "DELETE FROM service WHERE id = ?"
	res, err := ss.db.ExecContext(ctx, query, id)
	if err != nil {

		ss.logger.Warn(
			"Failed to remove service",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("service id", id),
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
			slog.Int("service id", id),
			"err", err,
		)
		return err
	}

	if affectedRows == 0 {

		ss.logger.Info(
			"Service does not exist",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("service id", id),
		)
		return repoerrors.ErrNotFound
	}

	ss.logger.Info(
		"Service successfully removed",
		slog.String("op", op),
		slog.Int("service id", id),
	)
	return nil
}

// GetList retrieves a list of services from the database.
//
// ctx - the context for the operation.
// *core.List[*core.Service], error - returns a list of services and an error, if any.
func (ss *ServiceSQL) GetList(ctx context.Context) (*core.List[*core.Service], error) {
	const op = "ServiceSQL.GetList"

	query := "SELECT id, name FROM service"
	rows, err := ss.db.QueryContext(ctx, query)
	if err != nil {
		ss.logger.Warn(
			"Failed to get service list",
			slog.String("op", op),
			slog.String("query", query),
			"err", err,
		)
		return nil, err
	}

	serviceList := make(core.List[*core.Service], 0)

	for rows.Next() {
		var (
			id   int
			name string
		)
		err = rows.Scan(&id, &name)
		if err != nil {
			ss.logger.Warn(
				"Failed to scan service row",
				slog.String("op", op),
				slog.String("query", query),
				"err", err,
			)
			return nil, err
		}

		service := core.NewService(id, name)
		serviceList[id] = &service
	}

	ss.logger.Info(
		"Service list successfully received",
		slog.String("op", op),
		slog.Int("service count", len(serviceList)),
	)
	return &serviceList, nil
}

// Update updates a service in the database.
//
// ctx: the context in which the update operation should be performed.
// s: the service object to be updated.
// error: an error if the update operation fails.
func (ss *ServiceSQL) Update(ctx context.Context, s *core.Service) error {
	const op = "ServiceSQL.Update"

	query := "UPDATE service SET name = ? WHERE id = ?"
	_, err := ss.db.ExecContext(ctx, query, s.Name(), s.Id())
	if err != nil {

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ss.logger.Info(
				"Service already exists",
				slog.String("op", op),
				slog.String("query", query),
				slog.Int("service id", s.Id()),
			)
			return repoerrors.ErrAlreadyExists
		}

		ss.logger.Warn(
			"Failed to update service",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("service id", s.Id()),
			"err", err,
		)
		return err
	}

	ss.logger.Info(
		"Service successfully updated",
		slog.String("op", op),
		slog.Int("service id", s.Id()),
		slog.String("service name", s.Name()),
	)
	return nil
}
