package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	
	return pool, nil
}
