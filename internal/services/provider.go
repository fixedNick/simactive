package services

import (
	"context"
	"simactive/internal/core"
)

type ProviderService struct {
	repository Repository[*core.Provider]
}

func NewProviderService(repo Repository[*core.Provider]) *ProviderService {
	ss := &ProviderService{
		repository: repo,
	}
	return ss
}

func (ps *ProviderService) GetProviderList(ctx context.Context) (*core.List[*core.Provider], error) {
	panic("implement")
}
