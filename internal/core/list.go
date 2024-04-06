package core

import (
	"database/sql"
	"fmt"
)

type Scannable interface {
	Scan(rows *sql.Rows) (Scannable, int, error)
}

type List[T Scannable] map[int]T

func NewSimList[T Scannable]() List[T] {
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
