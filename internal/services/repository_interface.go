package services

import (
	"context"
	"simactive/internal/core"
)

type Repository[T core.Scannable] interface {
	Save(ctx context.Context, obj core.Scannable) (id int, err error)
	GetList(ctx context.Context) (list *core.List[T], err error)
	Remove(ctx context.Context, idx int) (err error)
	ByID(ctx context.Context, id int) (T, error)
	Update(ctx context.Context, s *T) (err error)
}