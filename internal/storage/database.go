package storage

import (
	"accessCloude/internal/config"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Database struct {
	PConn *pgxpool.Pool
	MConn *minio.Client
	Salt  string
}

func NewDatabase(cfg *config.Config) *Database {
	ctx := context.Background()

	pconn := NewPostgresConn(ctx, cfg)
	mconn := NewMinioConn(ctx, cfg)

	return &Database{
		PConn: pconn,
		MConn: mconn,
	}
}

func NewPostgresConn(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Database)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse configuration: %v", err)
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil
	}

	return pool
}

func NewMinioConn(ctx context.Context, cfg *config.Config) *minio.Client {

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.SSL,
	})
	if err != nil {
		log.Fatalln("Unable to initialize Minio Client: ", err)
		return nil
	}

	return minioClient
}
