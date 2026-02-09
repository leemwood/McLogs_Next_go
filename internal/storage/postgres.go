package storage

import (
	"context"
	"fmt"
	"mclogs-go/internal/config"
	"mclogs-go/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
	cfg  *config.StorageConfig
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	// 1. Connect to default 'postgres' database to ensure 'mclogs' exists
	defaultDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable",
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	defaultPool, err := pgxpool.New(ctx, defaultDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to default postgres db: %w", err)
	}
	defer defaultPool.Close()

	// Check if database exists
	var exists bool
	err = defaultPool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.Database.Postgres.DB).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		_, err = defaultPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Postgres.DB))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
	}

	// 2. Connect to the target database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.DB,
	)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to target database: %w", err)
	}

	// Initialize table if not exists
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS logs (
			id VARCHAR(255) PRIMARY KEY,
			content TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP WITH TIME ZONE
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs table: %w", err)
	}

	return &PostgresStorage{
		pool: pool,
		cfg:  &cfg.Storage,
	}, nil
}

func (s *PostgresStorage) Put(ctx context.Context, content string) (string, error) {
	rawID := GenerateRawID()
	id := GetFullID(s.cfg.CurrentID, rawID)
	expiresAt := time.Now().Add(time.Duration(s.cfg.TTL) * time.Second)

	_, err := s.pool.Exec(ctx, "INSERT INTO logs (id, content, expires_at) VALUES ($1, $2, $3)", id, content, expiresAt)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostgresStorage) Get(ctx context.Context, id string) (*models.Log, error) {
	var log models.Log
	err := s.pool.QueryRow(ctx, "SELECT id, content, created_at, expires_at FROM logs WHERE id = $1", id).Scan(
		&log.ID, &log.Content, &log.CreatedAt, &log.ExpiresAt,
	)
	if err != nil {
		// pgx returns err if no rows found
		return nil, nil 
	}
	return &log, nil
}

func (s *PostgresStorage) Delete(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM logs WHERE id = $1", id)
	return err
}

func (s *PostgresStorage) Renew(ctx context.Context, id string) error {
	expiresAt := time.Now().Add(time.Duration(s.cfg.TTL) * time.Second)
	_, err := s.pool.Exec(ctx, "UPDATE logs SET expires_at = $1 WHERE id = $2", expiresAt, id)
	return err
}
