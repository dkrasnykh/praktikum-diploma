package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
)

func New(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnectTimeout)*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres connection error: %w", err)
	}
	if err = migrate(pool, 1); err != nil {
		return nil, fmt.Errorf("postgres migration error: %w", err)
	}
	return pool, nil
}
