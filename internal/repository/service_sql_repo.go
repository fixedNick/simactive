package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"simactive/internal/core"
)

type ServiceSqlRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewServiceSQLRepository(db *sql.DB, logger *slog.Logger) ServiceSqlRepository {
	return ServiceSqlRepository{
		db:     db,
		logger: logger,
	}
}

func (sr ServiceSqlRepository) Save(ctx context.Context, service core.Service) (id int, err error) {
	panic("implement")
}
func (sr ServiceSqlRepository) GetServiceList(ctx context.Context) (list *core.ServiceList, err error) {
	panic("implement")
}
func (sr ServiceSqlRepository) Remove(ctx context.Context, id int) (err error) {
	panic("implement")
}
func (sr ServiceSqlRepository) Update(ctx context.Context, s *core.Service) (err error) {
	panic("implement")
}
