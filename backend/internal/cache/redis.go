package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/abhinavmaity/linear-lite/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg config.Config) (*redis.Client, error) {
	options, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
