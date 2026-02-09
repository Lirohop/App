package main

import (
	"app/internal/config"
	"app/internal/database"
	"app/internal/handler"
	"app/internal/repository"
	"app/internal/service"
	 httpSwagger "github.com/swaggo/http-swagger"
	 _ "app/docs"
	"fmt"
	"log/slog"
	"net/http"
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

	pool, err := database.NewDatabase(cfg, logger)
	if err != nil {
		logger.Error("Failed to connect to database, exiting", "error", err)
		panic(err)
	}
	defer pool.Close()

	logger.Debug("Database connection object created, passing to repository")

	rep := repository.NewSubscriptionRepository(pool, logger)
	logger.Info("Repository initialized")

	subService := service.NewSubscriptionService(rep, logger)

	subHandler := handler.NewSubscriptionHandler(subService, logger)

	http.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			subHandler.Create(w, r)
		case http.MethodGet:
			subHandler.List(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/subscriptions/get", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		subHandler.GetByID(w, r)
	})

	http.HandleFunc("/subscriptions/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		subHandler.Delete(w, r)
	})

	http.HandleFunc("/subscriptions/total-cost", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		subHandler.TotalCost(w, r)
	})

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	logger.Debug("Startup complete, ready to handle requests")

	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info("HTTP server listening", "addr", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("Server stopped unexpectedly", "error", err)
	}
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
