package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// Create an empty context for the database connection. This will allow the connection to be established without any specific timeout or cancellation.
	pool, err := pgxpool.New(
		context.Background(),
		databaseUrl,
	)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}

	fmt.Println("Connected to PostgreSQL")

	return pool, nil
}
