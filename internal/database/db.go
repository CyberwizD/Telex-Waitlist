package database

import (
	"context"
	"time"

	"github.com/CyberwizD/Telex-Waitlist/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect initializes a pgx connection pool.
func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = 5
	poolConfig.MinConns = 1
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 15 * time.Minute

	return pgxpool.NewWithConfig(ctx, poolConfig)
}
