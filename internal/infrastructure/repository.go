package repository

import (
	"database/sql"
	"log/slog"
	providerrepository "simactive/internal/infrastructure/provider"
	servicerepository "simactive/internal/infrastructure/service"
	simrepository "simactive/internal/infrastructure/sim"
	usedrepository "simactive/internal/infrastructure/used"
)

type Repository struct {
	SimRepository      *simrepository.SimRepository
	ServiceRepository  *servicerepository.ServiceRepository
	ProviderRepository *providerrepository.ProviderRepository
	UsedRepository     *usedrepository.UsedRepository
}

func NewRepository(logger *slog.Logger, db *sql.DB) *Repository {
	return &Repository{
		SimRepository: simrepository.NewSimRepository(
			logger,
			db,
			simrepository.NewSimInMemoryRepository(logger),
			simrepository.NewSimSQLRepository(db, logger),
		),

		ServiceRepository: servicerepository.NewServiceRepository(
			logger,
			db,
			servicerepository.NewServiceInMemoryRepository(logger),
			servicerepository.NewServiceSQLRepository(db, logger),
		),

		ProviderRepository: providerrepository.NewProviderRepository(
			logger,
			db,
			providerrepository.NewProviderInMemory(logger),
			providerrepository.NewProviderSQL(db, logger),
		),
		UsedRepository: usedrepository.NewUsedRepository(
			logger,
			db,
			usedrepository.NewUsedInMemoryRepository(logger),
			usedrepository.NewUsedSQLRepository(db, logger),
		),
	}
}
