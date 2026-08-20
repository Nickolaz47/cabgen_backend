package services

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskEnqueuer interface {
	EnqueueContext(ctx context.Context, task *asynq.Task,
		opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type TaskCanceller interface {
	CancelProcessing(id string) error
}
