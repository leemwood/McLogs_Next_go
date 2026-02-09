package storage

import (
	"context"
	"mclogs-go/internal/models"
)

type Storage interface {
	Put(ctx context.Context, content string) (string, error)
	Get(ctx context.Context, id string) (*models.Log, error)
	Delete(ctx context.Context, id string) error
	Renew(ctx context.Context, id string) error
}
