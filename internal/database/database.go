package database

import (
	. "app/internal/config"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func NewDatabase(cfg *Config, logger *slog.Logger) (*pgx.Conn, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	logger.Debug("Creating database connection", "dsn", dsn)

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		logger.Error("Database ping failed", "error", err)
		return nil, err
	}

	logger.Info("Database connection established successfully")
	return conn, nil
}

