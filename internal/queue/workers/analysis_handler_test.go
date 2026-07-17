package workers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/workers"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func TestAnalysisTaskHandlerProcessTask(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		analysisID := uuid.New()
		mockService := &mocks.MockAnalysisRunnerService{
			RunFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				assert.Equal(t, analysisID, receivedID)
				return nil
			},
		}
		handler := workers.NewAnalysisTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AnalysisProcessPayload{
			AnalysisID: analysisID})
		task := asynq.NewTask(tasks.TaskTypeAnalysisProcess, payloadBytes)

		err := handler.ProcessTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("Error - JSON Unmarshal", func(t *testing.T) {
		mockService := &mocks.MockAnalysisRunnerService{}
		handler := workers.NewAnalysisTaskHandler(mockService)

		task := asynq.NewTask(tasks.TaskTypeAnalysisProcess,
			[]byte(`{"analysis_id": "invalid-uuid"`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, asynq.SkipRetry)
		assert.ErrorContains(t, err, "json unmarshal failed:")
	})

	t.Run("Error - Service Failure", func(t *testing.T) {
		analysisID := uuid.New()
		mockService := &mocks.MockAnalysisRunnerService{
			RunFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				return errors.New("pipeline crashed")
			},
		}
		handler := workers.NewAnalysisTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AnalysisProcessPayload{
			AnalysisID: analysisID})
		task := asynq.NewTask(tasks.TaskTypeAnalysisProcess, payloadBytes)

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.EqualError(t, err, "pipeline crashed")
	})

	t.Run("Error - Unknown Task Type", func(t *testing.T) {
		mockService := &mocks.MockAnalysisRunnerService{}
		handler := workers.NewAnalysisTaskHandler(mockService)

		task := asynq.NewTask("analysis:alien_task", []byte(`{}`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "unknown task type: analysis:alien_task")
	})
}