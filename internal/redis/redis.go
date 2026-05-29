package redis

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	opt, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Redis URL: %v", err)
	}

	opt.PoolSize = 120
	opt.DB = 0
	opt.MaxRetries = 3

	client := redis.NewClient(opt)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("Failed to connect to Redis: %v", err)
	}

	return client, nil
}
