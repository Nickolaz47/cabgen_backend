package mocks

import (
	"context"

	"github.com/hibiken/asynq"
)

type MockTaskEnqueuer struct {
	EnqueueContextFunc func(ctx context.Context, task *asynq.Task,
		opts ...asynq.Option) (*asynq.TaskInfo, error)
}

func (m *MockTaskEnqueuer) EnqueueContext(ctx context.Context,
	task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	if m.EnqueueContextFunc != nil {
		return m.EnqueueContextFunc(ctx, task, opts...)
	}

	return &asynq.TaskInfo{ID: "mock-task-id", Queue: "emails"}, nil
}

type MockTaskCanceller struct {
	CancelProcessingFunc func(id string) error
	Called               bool
	CalledWith           string
	CallCount            int
}

func (m *MockTaskCanceller) CancelProcessing(id string) error {
	m.Called = true
	m.CalledWith = id
	m.CallCount++
	if m.CancelProcessingFunc != nil {
		return m.CancelProcessingFunc(id)
	}
	return nil
}
