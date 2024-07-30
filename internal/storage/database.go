package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Conn *pgxpool.Pool
	Salt string
}

func NewDatabase(url string) *Database {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("Unable to parse configuration: %v", err)
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil
	}

	return &Database{
		Conn: pool,
	}
}
