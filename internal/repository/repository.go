package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"simactive/internal/core"
)

type Sql struct {
	Db *sql.DB
}

type Repository[T core.Scannable] struct {
	Sql

	InMemoryList core.List[T]
	logger       *slog.Logger
}

func NewRepository[T core.Scannable](ctx context.Context, db *sql.DB, logger *slog.Logger) *Repository[T] {
	repository := &Repository[T]{
		logger: logger,
	}
	repository.Db = db

	// Load list from sql
	// TODO:
	// Force load from sql, now it loads it cuz of local list is empty and ONLY
	if list, err := repository.GetList(ctx); err == nil {
		repository.InMemoryList = *list
	}
	return repository
}

// obj allowd types: [core.Sim | core.Service | core.Provider | core.Used]
func (r *Repository[T]) Save(ctx context.Context, obj core.Scannable) (int, error) {

	var query string
	var args []interface{}

	switch obj := any(obj).(type) {
	case core.Sim:

		// Is Local Storage Contains it ?
		if _, ok := r.InMemoryList[obj.Id()]; ok {
			return 0, fmt.Errorf("sim already in the in-memory list")
		}

		// Prepare SQL query
		query = "INSERT INTO sim VALUES (?, ?, ?, ?, ?, ?)"
		args = []interface{}{
			0, obj.Number(), obj.ProviderID(), obj.IsActivated(), obj.ActivateUntil(), obj.IsBlocked(),
		}
	case core.Service:

		// Is Llocal Storage Contains it ?
		if _, ok := r.InMemoryList[obj.Id()]; ok {
			return 0, fmt.Errorf("service already in the in-memory list")
		}

		// Prepare SQL query
		query = "INSER INTO service VALUES (?, ?)"
		args = []interface{}{
			0, obj.Name(),
		}
	case core.Provider:
		panic("implement")
	default:
		panic("implement")
	}

	res, err := r.Sql.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	return int(id), err
}

func (r *Repository[T]) Remove(ctx context.Context, id int) error {
	panic("implement")
}
func (r *Repository[T]) GetList(ctx context.Context) (*core.List[T], error) {

	if len(r.InMemoryList) != 0 {
		return &r.InMemoryList, nil
	}

	var v T
	var query string

	switch interface{}(v).(type) {
	case core.Sim:
		query = "SELECT * FROM sim"
	case core.Service:
		query = "SELECT * FROM service"
	default:
		panic("impossible")
	}

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	list := make(core.List[T])

	for {
		if !rows.Next() {
			break
		}

		v, key, err := v.Scan(rows)
		if err != nil {
			return nil, err
		}
		list[key] = v.(T)
	}
	return &list, nil
}
func (r *Repository[T]) ByID(ctx context.Context, id int) (T, error) {
	panic("implement")
}
func (r *Repository[T]) Update(ctx context.Context, s *T) error {
	panic("implement")
}
