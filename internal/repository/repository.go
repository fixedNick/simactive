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

type Repository[T core.DBModel] struct {
	Sql

	InMemoryList core.List[T]
	logger       *slog.Logger
}

func NewRepository[T core.DBModel](ctx context.Context, db *sql.DB, logger *slog.Logger) *Repository[T] {
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

func (r *Repository[T]) Save(ctx context.Context, obj core.DBModel) (int, error) {

	var query string
	var args []interface{}

	// setup query & args
	switch obj := any(obj).(type) {
	case *core.Sim:
		// Prepare SQL query
		query = "INSERT INTO sim (id, number, provider_id, is_activated, activate_until, is_blocked) VALUES (?, ?, ?, ?, ?, ?)"
		args = []interface{}{
			0, obj.Number(), obj.ProviderID(), obj.IsActivated(), obj.ActivateUntil(), obj.IsBlocked(),
		}
	case *core.Service:
		// Prepare SQL query
		query = "INSER INTO service VALUES (?, ?)"
		args = []interface{}{
			0, obj.Name(),
		}
	case *core.Provider:
		panic("implement")
	default:
		panic("implement")
	}

	res, err := r.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	fmt.Println("Inserted id: ", id)

	// Save into in-memory storage
	obj.SetKey(int(id))
	r.InMemoryList[int(id)] = obj.(T)
	return int(id), err
}
func (r *Repository[T]) Remove(ctx context.Context, id int) error {
	// first check is that exist locally
	obj, ok := r.InMemoryList[id]
	if !ok {
		// if not - return error
		return fmt.Errorf("locally not found object of type `%T` with id `%d`", obj, id)
	}

	// then delete locally
	delete(r.InMemoryList, id)

	// then delete from sql
	var query string

	switch any(obj).(type) {
	case *core.Sim:
		query = "DELETE FROM sim WHERE id = ?"
	case *core.Service:
		query = "DELTE FROM service WHERE id = ?"
	}

	_, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error ocured on deleting row from sql. Delete type: `%T`, id: `%d`", obj, id)
	}
	return nil
}
func (r *Repository[T]) GetList(ctx context.Context) (*core.List[T], error) {
	if len(r.InMemoryList) != 0 {
		return &r.InMemoryList, nil
	}

	var v T
	var query string

	switch interface{}(v).(type) {
	case *core.Sim:
		query = "SELECT * FROM sim"
	case *core.Service:
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

		key, err := v.ScanRows(rows)
		if err != nil {
			return nil, err
		}
		list[key] = v
	}
	return &list, nil
}
func (r *Repository[T]) ByID(ctx context.Context, id int) (T, error) {

	if obj, ok := r.InMemoryList[id]; ok {
		return obj, nil
	}

	var v T
	var query string
	var args []interface{}

	switch any(v).(type) {
	case *core.Sim:
		query = "SELECT * FROM sim WHERE id = ?"
		args = []interface{}{
			id,
		}
	case *core.Service:
		query = "SELECT * FROM service WHERE id = ?"
		args = []interface{}{
			id,
		}
	case *core.Provider:
	case *core.Used:
	}

	row := r.Db.QueryRowContext(ctx, query, args...)
	err := v.ScanRow(row)
	if err != nil {
		return v, err
	}

	return v, nil
}
func (r *Repository[T]) Update(ctx context.Context, s T) error {
	// check is object exist locally
	// if not - return err
	// send update query

	if _, ok := r.InMemoryList[s.GetKey()]; !ok {
		return fmt.Errorf("sim: %v. Not found in memory", s)
	}
	// Save local
	r.InMemoryList[s.GetKey()] = s

	var query string
	var args []interface{}

	switch obj := any(s).(type) {
	case *core.Sim:
		query = "UPDATE sim SET number = ?, provider_id = ?, is_activated = ?, activate_until = ?, is_blocked = ? WHERE id = ?"
		args = []interface{}{
			obj.Number(),
			obj.ProviderID(),
			obj.IsActivated(),
			obj.ActivateUntil(),
			obj.IsBlocked(),
			obj.Id(),
		}
	case *core.Service:
		query = "UPDATE service SET name = ? WHERE id = ?"
		args = []interface{}{
			obj.Name(),
			obj.Id(),
		}
	case *core.Provider:
	case *core.Used:
	}

	_, err := r.Db.ExecContext(ctx, query, args...)
	return err
}
