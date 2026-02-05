package main

import (
	"app/internal/config"
	"app/internal/database"
	"app/internal/repository"
	"context"
	"log/slog"
	"os"
)

const (
	logDebug = "debug"
	logDev   = "dev" //default
	logProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	logger := setupLogger(cfg.App.LogLevel)
	slog.SetDefault(logger)
	logger.Info("application starter", "port", cfg.App.Port, "log_level", cfg.App.LogLevel)

	conn, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database, exiting", "error", err)
		panic(err)
	}
	defer conn.Close(context.Background())

	logger.Debug("Database connection object created, passing to repository")

	rep := repository.NewSubscriptionRepository(conn, logger)
	logger.Info("Repository initialized")

	//TODO regist handlers

	logger.Debug("Startup complete, ready to handle requests")

}

func setupLogger(logLevel string) (logger *slog.Logger) {

	switch logLevel {
	case logDebug:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case logDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case logProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return
}
