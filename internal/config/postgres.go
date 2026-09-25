package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	_ = godotenv.Load()
	DBDRIVER := os.Getenv("DB_DRIVER")
	switch DBDRIVER {
	case "local":
		os.Setenv("DATABASE_URL", os.Getenv("LOCAL_DATABASE_URL"))
	case "neon":
		os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL"))
	default:
		return nil, fmt.Errorf(
			"unsupported DB_DRIVER: %q",
			DBDRIVER,
		)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment is not found or not set")
	}

	// Create an empty context for the database connection. This will allow the connection to be established without any specific timeout or cancellation.
	pool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		pool.Close()
		return nil, err
	}

	if(DBDRIVER == "neon") {
		fmt.Println("Connected to Neon PostgreSQL")
	} else {
		fmt.Println("Connected to Local PostgreSQL")
	}

	return pool, nil
}
