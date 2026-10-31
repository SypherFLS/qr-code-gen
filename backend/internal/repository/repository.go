package repository

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, id string, data []byte, ttl time.Duration) error
	Get(ctx context.Context, id string) ([]byte, error)
}
