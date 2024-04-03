package main

import (
	"database/sql"
	"log/slog"
	"simactive/internal/config"
	"simactive/internal/core/grpc"
	"simactive/internal/repository"
	"simactive/internal/services"
	coresql "simactive/internal/sql"
)

func main() {

	// Initialize config object
	cfg := config.MustLoad()

	// Initialize logger
	logger := slog.Default()

	// Init db
	db := coresql.MustInit()

	// Init services
	simService := InitServices(db, logger)

	// Init gRPC Server
	gs := grpc.NewGRPCServer(cfg)
	// Run gRPC server
	gs.MustRun(simService)

	// gracefull shutdown
	//...
}

func InitServices(db *sql.DB, logger *slog.Logger) *services.SimService {
	simService := services.NewSimService(
		repository.NewInMemoryRepository(),
		repository.NewSQLRepository(db, logger),
	)

	return simService
}
