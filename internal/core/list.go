package core

import (
	"database/sql"
	"simactive/internal/infrastructure/repoerrors"
)

type DBModel interface {
	Scannable
	Keyable
}

type Keyable interface {
	SetKey(id int)
	GetKey() int
}

type Scannable interface {
	ScanRows(rows *sql.Rows) (int, error)
	ScanRow(rows *sql.Row) error
}

type List[T DBModel] map[int]T

// NewSimList creates a new List of type T where T implements DBModel interface
func NewSimList[T DBModel]() List[T] {
	return make(List[T])
}

func (l *List[T]) ByID(id int) (T, error) {
	obj, ok := (*l)[id]
	if !ok {
		return obj, repoerrors.ErrNotFound
	}
	return obj, nil
}

func (l *List[T]) ContainsFunc(cf func(T) bool) (T, bool) {
	for _, v := range *l {
		if cf(v) {
			return v, true
		}
	}

	return *new(T), false
}
