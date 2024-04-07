package services

import (
	"context"
	"simactive/internal/core"
)

type UsedService struct {
	repository Repository[*core.Used]
}

func NewUsedService(repo Repository[*core.Used]) *UsedService {
	ss := &UsedService{
		repository: repo,
	}
	return ss
}

func (us *UsedService) UseSimForService(ctx context.Context, simId int, serviceId int) {
	panic("implement")
}
