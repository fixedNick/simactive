package grpc

import (
	"context"
	"simactive/internal/core"
)

type SimService interface {
	Add(ctx context.Context, s core.Sim) error
	GetByID(ctx context.Context, id int) (sim core.Sim, err error)
	Remove(ctx context.Context, s core.Sim) error
}
