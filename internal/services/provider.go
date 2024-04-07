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

// GetProviderList retrieves a list of providers.
// Returns a list of Provider objects and an error.
func (ps *ProviderService) GetProviderList(ctx context.Context) (*core.List[*core.Provider], error) {
	return ps.repository.GetList(ctx)
}

// Add adds a Provider to the ProviderService.
// ctx - the context in which the operation is performed.
// p - the Provider to be added.
// Returns an int represents the ID of the added Provider and an error.
func (ps *ProviderService) Add(ctx context.Context, p *core.Provider) (int, error) {
	return ps.repository.Save(ctx, p)
}
