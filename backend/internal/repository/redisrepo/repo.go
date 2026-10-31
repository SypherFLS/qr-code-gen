package redisrepo

import (
	"context"
	"errors"
	"fmt"
	"qrgen/internal/apperrors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repo struct {
	rdb *redis.Client
}

func NewRepo(rdb *redis.Client) *Repo {
	return &Repo{
		rdb: rdb,
	}
}
func (r *Repo) Save(ctx context.Context, id string, data []byte, ttl time.Duration) error {
	key := "qr:" + id

	if err := r.rdb.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("save qr code: %w", err)
	}

	return nil
}

func (r *Repo) Get(ctx context.Context, id string) ([]byte, error) {
	key := "qr:" + id

	data, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("get qr code: %w", err)
	}

	return data, nil
}
