package mocks

import (
	"context"

	"github.com/hibiken/asynq"
)

// MockTaskQueue implements services.TaskQueue with per-method call
// tracking, so tests can assert on enqueue, delete and cancel separately.
type MockTaskQueue struct {
	EnqueueContextFunc func(ctx context.Context, task *asynq.Task,
		opts ...asynq.Option) (*asynq.TaskInfo, error)
	DeleteTaskFunc       func(queue, id string) error
	CancelProcessingFunc func(id string) error

	EnqueueCalled bool
	EnqueueCount  int

	DeleteCalled    bool
	DeleteWith      string
	DeleteQueue     string
	DeleteCallCount int

	CancelCalled bool
	CancelWith   string
	CancelCount  int
}

func (m *MockTaskQueue) EnqueueContext(ctx context.Context,
	task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	m.EnqueueCalled = true
	m.EnqueueCount++
	if m.EnqueueContextFunc != nil {
		return m.EnqueueContextFunc(ctx, task, opts...)
	}

	return &asynq.TaskInfo{ID: "mock-task-id", Queue: "emails"}, nil
}

func (m *MockTaskQueue) DeleteTask(queue, id string) error {
	m.DeleteCalled = true
	m.DeleteWith = id
	m.DeleteQueue = queue
	m.DeleteCallCount++
	if m.DeleteTaskFunc != nil {
		return m.DeleteTaskFunc(queue, id)
	}
	return nil
}

func (m *MockTaskQueue) CancelProcessing(id string) error {
	m.CancelCalled = true
	m.CancelWith = id
	m.CancelCount++
	if m.CancelProcessingFunc != nil {
		return m.CancelProcessingFunc(id)
	}
	return nil
}
