package queue

import (
	"errors"

	"github.com/hibiken/asynq"
)

func NewAsynqClient(redisAddr string) (*asynq.Client, error) {
	if redisAddr == "" {
		return nil, errors.New("Redis address is empty")
	}

	opt := asynq.RedisClientOpt{
		Addr:     redisAddr,
		PoolSize: 120,
		DB:       0,
	}

	client := asynq.NewClient(opt)
	if err := client.Ping(); err != nil {
		return nil, err
	}

	return client, nil
}
