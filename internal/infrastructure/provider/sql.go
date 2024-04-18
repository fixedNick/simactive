package providerrepository

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

type ProviderSQL struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewProviderSQL(db *sql.DB, logger *slog.Logger) *ProviderSQL {
	return &ProviderSQL{
		db:     db,
		logger: logger,
	}
}

// Add adds a new provider with the given name to the database.
//
// ctx is the context for the operation.
// name is the name of the provider to add.
// Returns the ID of the newly added provider and any error encountered.
func (ps *ProviderSQL) Add(ctx context.Context, name string) (int, error) {
	const op = "ProviderSQL.Add"

	query := "INSERT INTO provider (name) VALUES (?)"
	res, err := ps.db.ExecContext(ctx, query, name)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ps.logger.Info(
				"Provider already exists",
				slog.String("op", op),
				slog.String("query", query),
				slog.String("provider name", name),
			)
			return 0, repoerrors.ErrAlreadyExists
		}
		ps.logger.Warn(
			"Failed to add provider",
			slog.String("op", op),
			slog.String("query", query),
			slog.String("provider name", name),
			sl.Err(err),
		)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		ps.logger.Warn(
			"Failed to receive last insert id after query",
			slog.String("op", op),
			slog.String("query", query),
			slog.String("provider name", name),
			sl.Err(err),
		)
		return 0, err
	}

	return int(id), nil
}
func (ps *ProviderSQL) GetList(ctx context.Context) (*core.List[*core.Provider], error) {
	const op = "ProviderSQL.GetList"

	query := "SELECT id, name FROM provider"
	rows, err := ps.db.QueryContext(ctx, query)
	if err != nil {
		ps.logger.Warn(
			"Failed to get provider list",
			slog.String("op", op),
			slog.String("query", query),
			sl.Err(err),
		)
		return nil, err
	}

	providerList := make(core.List[*core.Provider], 0)
	for rows.Next() {
		var (
			id   int
			name string
		)
		err = rows.Scan(&id, &name)
		if err != nil {
			ps.logger.Warn(
				"Failed to scan provider row",
				slog.String("op", op),
				slog.String("query", query),
				sl.Err(err),
			)
			return nil, err
		}

		p := core.NewProvider(id, name)
		providerList[id] = &p
	}

	ps.logger.Info(
		"Provider list successfully retrieved",
		slog.String("op", op),
		slog.Int("provider count", len(providerList)),
	)
	return &providerList, nil
}
func (ps *ProviderSQL) ByID(ctx context.Context, id int) (*core.Provider, error) {
	const op = "ProviderSQL.ByID"

	query := "SELECT name FROM provider WHERE id = ?"

	var name string
	err := ps.db.QueryRowContext(ctx, query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ps.logger.Info(
				"Provider does not exist",
				slog.String("op", op),
				slog.String("query", query),
				slog.Int("provider id", id),
			)
			return nil, repoerrors.ErrNotFound
		}

		ps.logger.Warn(
			"Failed to get provider",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("provider id", id),
			sl.Err(err),
		)
		return nil, err
	}

	p := core.NewProvider(id, name)

	ps.logger.Info(
		"Provider successfully retrieved",
		slog.String("op", op),
		slog.Int("provider id", id),
		slog.String("provider name", name),
		slog.Any("provider", p),
	)
	return &p, nil
}

// ByName retrieves a provider by name from the database.
//
// ctx: the context for the operation.
// name: the name of the provider to retrieve.
// *core.Provider, error: returns the provider information and any errors encountered.
func (ps *ProviderSQL) ByName(ctx context.Context, name string) (*core.Provider, error) {
	const op = "ProviderSQL.ByName"

	query := "SELECT id FROM provider WHERE name = ?"

	var id int
	err := ps.db.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ps.logger.Info(
				"Provider does not exist",
				slog.String("op", op),
				slog.String("query", query),
				slog.String("provider name", name),
			)
			return nil, repoerrors.ErrNotFound
		}

		ps.logger.Warn(
			"Failed to get provider",
			slog.String("op", op),
			slog.String("query", query),
			slog.String("provider name", name),
			sl.Err(err),
		)
		return nil, err
	}
	p := core.NewProvider(id, name)

	ps.logger.Info(
		"Provider successfully retrieved",
		slog.String("op", op),
		slog.Int("provider id", id),
		slog.String("provider name", name),
		slog.Any("provider", p),
	)
	return &p, nil
}

// Remove removes a provider from the database by ID.
//
// ctx: the context for the database operation.
// id: the ID of the provider to remove.
// error: returns an error if the operation fails.
func (ps *ProviderSQL) Remove(ctx context.Context, id int) error {
	const op = "ProviderSQL.Remove"

	query := "DELETE FROM provider WHERE id = ?"
	res, err := ps.db.ExecContext(ctx, query, id)
	if err != nil {
		ps.logger.Warn(
			"Failed to remove provider",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("provider id", id),
			sl.Err(err),
		)
		return err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		ps.logger.Warn(
			"Failed to receive affected rows after query",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("provider id", id),
			sl.Err(err),
		)
		return err
	}

	if affectedRows == 0 {
		ps.logger.Info(
			"Provider does not exist",
			slog.String("op", op),
			slog.String("query", query),
			slog.Int("provider id", id),
		)
		return repoerrors.ErrNotFound
	}

	ps.logger.Info(
		"Provider successfully removed",
		slog.String("op", op),
		slog.Int("provider id", id),
	)
	return nil
}
