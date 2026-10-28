package services

import (
	"context"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

// CommitGrace delays every task dispatch (asynq.ProcessIn) long enough for
// the DB write of the same request to commit
const CommitGrace = 3 * time.Second

func DeleteScheduledTask(ctx context.Context, tq TaskQueue, logger *zap.Logger,
	service, funcName, queue, taskID string) {
	if err := tq.DeleteTask(queue, taskID); err != nil {
		logger.Error("Service Error", logging.ServiceLogging(ctx,
			service, funcName, logging.TaskCancelError, err,
			zap.String("task_id", taskID),
		)...)
	}
}

type TaskQueue interface {
	EnqueueContext(ctx context.Context, task *asynq.Task,
		opts ...asynq.Option) (*asynq.TaskInfo, error)
	CancelProcessing(id string) error
	DeleteTask(queue, id string) error
}
