package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"simactive/internal/config"
	"simactive/internal/core/grpc"
	"simactive/internal/repository"
	"simactive/internal/services"
	coresql "simactive/internal/sql"
	"syscall"
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
	go func() {
		gs.MustRun(simService)
	}()

	// gracefull shutdown
	//...

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	gs.Stop()
	log.Print("Gracefull shutdown")
}

func InitServices(db *sql.DB, logger *slog.Logger) *services.SimService {
	simService := services.NewSimService(
		repository.NewSimInMemoryRepository(),
		repository.NewSimSQLRepository(db, logger),
	)
	return simService
}
