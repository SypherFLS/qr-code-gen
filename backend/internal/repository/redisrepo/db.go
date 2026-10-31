package redisrepo

import (
	"context"
	"fmt"
	"qrgen/internal/config"

	"github.com/redis/go-redis/v9"
)

func InitClient(ctx context.Context, cfg config.Config) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Database.Host,
		Password: cfg.Database.Password,
		DB:       cfg.Database.DB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return rdb, nil
}
