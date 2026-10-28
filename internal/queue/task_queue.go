package queue

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
)

type taskQueue struct {
	client    *asynq.Client
	inspector *asynq.Inspector
}

// NewTaskQueue bundles the asynq client (enqueue) and inspector
// (cancel/delete) behind one services.TaskQueue.
func NewTaskQueue(client *asynq.Client,
	inspector *asynq.Inspector) services.TaskQueue {
	return &taskQueue{client: client, inspector: inspector}
}

func (t *taskQueue) EnqueueContext(ctx context.Context, task *asynq.Task,
	opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return t.client.EnqueueContext(ctx, task, opts...)
}

func (t *taskQueue) CancelProcessing(id string) error {
	return t.inspector.CancelProcessing(id)
}

func (t *taskQueue) DeleteTask(queue, id string) error {
	return t.inspector.DeleteTask(queue, id)
}
