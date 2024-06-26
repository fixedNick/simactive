package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"simactive/internal/config"
	"simactive/internal/core/grpc"
	repository "simactive/internal/infrastructure"
	"simactive/internal/lib/logger/handlers/slogpretty"
	"simactive/internal/services"
	coresql "simactive/internal/sql"
	"syscall"
)

func main() {

	// Initialize config object
	cfg := config.MustLoad()

	// Initialize logger
	logger := setupLogger()

	// Init db
	db := coresql.MustInit()

	// Init services
	repo := repository.NewRepository(logger, db)
	simService, serviceService, providerService, usedService := initServices(db, logger, repo)

	// Init gRPC Server
	gs := grpc.NewGRPCServer(cfg)
	// Run gRPC server
	go func() {
		gs.MustRun(logger, simService, serviceService, providerService, usedService)
	}()

	// gracefull shutdown
	//...

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	gs.Stop()
	log.Print("Gracefull shutdown")
}

func initServices(db *sql.DB, logger *slog.Logger, repo *repository.Repository) (*services.SimService, *services.ServiceService, *services.ProviderService, *services.UsedService) {

	simService := services.NewSimService(repo)

	serviceService := services.NewServiceService(repo)

	providerService := services.NewProviderService(repo)

	usedService := services.NewUsedService(repo)

	return simService, serviceService, providerService, usedService
}

func setupLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
