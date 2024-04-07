package core

import (
	"database/sql"
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
