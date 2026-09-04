package database

import (
	"context"
	"fmt"
	"log"

	"github.com/hurbbiee/todo-list-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(
	ctx context.Context,
	cfg config.DatabaseConfig,
) (*pgxpool.Pool, error) {
	connectContext, cancel := context.WithTimeout(
		ctx,
		cfg.ConnectTimeout,
	)
	defer cancel()

	log.Printf(
		"database connect timeout=%s ",
		cfg.ConnectTimeout,
	)

	pool, err := pgxpool.New(connectContext, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(connectContext); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return pool, nil
}
