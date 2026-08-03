package database

import "github.com/jackc/pgx/v5/pgxpool"

func NewDatabase(pool *pgxpool.Pool) (*pgxpool.Pool, error) {
	return pool, nil
}