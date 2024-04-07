package core

import (
	"database/sql"
	"fmt"
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

func NewSimList[T DBModel]() List[T] {
	list := make(List[T])
	return list
}

func (sl *List[T]) PtrByKey(key int) (*T, error) {
	s, ok := (*sl)[key]
	if !ok {
		return nil, fmt.Errorf("List [%v] does not contain key %v(%T)", *sl, key, key)
	}
	return &s, nil
}
