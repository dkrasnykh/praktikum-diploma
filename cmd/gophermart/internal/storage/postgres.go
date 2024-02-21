package storage

import (
	"context"
	"fmt"
	
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dkrasnykh/praktikum-diploma/cmd/gophermart/internal/config"
)

func New(cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("postgres connection error: %w", err)
	}
	if err = migrate(pool, 1); err != nil {
		return nil, fmt.Errorf("postgres migration error: %w", err)
	}
	return pool, nil
}
