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

func TestEmailTaskHandlerProcessTask(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Admin Alert Email", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendAdminAlertEmailFunc: func(ctx context.Context,
				newUserID uuid.UUID) error {
				assert.Equal(t, userID, newUserID)
				return nil
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AdminAlertEmailPayload{
			NewUserID: userID})
		task := asynq.NewTask(tasks.TaskTypeAdminAlertEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("Error - Admin Alert JSON Unmarshal", func(t *testing.T) {
		mockService := &mocks.MockEmailService{}
		handler := workers.NewEmailTaskHandler(mockService)

		task := asynq.NewTask(tasks.TaskTypeAdminAlertEmail, []byte(
			`{"new_user_id": "invalid-uuid"`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, asynq.SkipRetry)
		assert.ErrorContains(t, err, "json unmarshal failed:")
	})

	t.Run("Error - Admin Alert Service Failure", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendAdminAlertEmailFunc: func(ctx context.Context,
				newUserID uuid.UUID) error {
				return errors.New("service internal error")
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AdminAlertEmailPayload{
			NewUserID: userID})
		task := asynq.NewTask(tasks.TaskTypeAdminAlertEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.EqualError(t, err, "service internal error")
	})

	t.Run("Success - Welcome Email", func(t *testing.T) {
		userID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendWelcomeEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				assert.Equal(t, userID, receivedID)
				return nil
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.WelcomeEmailPayload{
			UserID: userID})
		task := asynq.NewTask(tasks.TaskTypeWelcomeEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("Error - Welcome Email JSON Unmarshal", func(t *testing.T) {
		mockService := &mocks.MockEmailService{}
		handler := workers.NewEmailTaskHandler(mockService)

		task := asynq.NewTask(tasks.TaskTypeWelcomeEmail, []byte(
			`broken json`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, asynq.SkipRetry)
	})

	t.Run("Success - Analysis Done Email", func(t *testing.T) {
		analysisID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendAnalysisDoneEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				assert.Equal(t, analysisID, receivedID)
				return nil
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AnalysisDoneEmailPayload{
			AnalysisID: analysisID})
		task := asynq.NewTask(tasks.TaskTypeAnalysisDoneEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("Error - Analysis Done Email JSON Unmarshal", func(t *testing.T) {
		mockService := &mocks.MockEmailService{}
		handler := workers.NewEmailTaskHandler(mockService)

		task := asynq.NewTask(tasks.TaskTypeAnalysisDoneEmail, []byte(
			`[1, 2, 3]`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, asynq.SkipRetry)
	})

	t.Run("Error - Analysis Done Service Failure", func(t *testing.T) {
		analysisID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendAnalysisDoneEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				return errors.New("smtp timeout")
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AnalysisDoneEmailPayload{
			AnalysisID: analysisID})
		task := asynq.NewTask(tasks.TaskTypeAnalysisDoneEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.EqualError(t, err, "smtp timeout")
	})

	t.Run("Error - Unknown Task Type", func(t *testing.T) {
		mockService := &mocks.MockEmailService{}
		handler := workers.NewEmailTaskHandler(mockService)

		task := asynq.NewTask("email:alien_task", []byte(`{}`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "unknown task type: email:alien_task")
	})

	t.Run("Success - Admin Ticket Email", func(t *testing.T) {
		ticketID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendAdminTicketEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				assert.Equal(t, ticketID, receivedID)
				return nil
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AdminTicketEmailPayload{
			TicketID: ticketID})
		task := asynq.NewTask(tasks.TaskTypeAdminTicketEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("Error - Admin Ticket Email JSON Unmarshal", func(t *testing.T) {
		mockService := &mocks.MockEmailService{}
		handler := workers.NewEmailTaskHandler(mockService)

		task := asynq.NewTask(tasks.TaskTypeAdminTicketEmail,
			[]byte(`broken json`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, asynq.SkipRetry)
		assert.ErrorContains(t, err, "json unmarshal failed:")
	})

	t.Run("Error - Admin Ticket Email Service Failure", func(t *testing.T) {
		ticketID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendAdminTicketEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				return errors.New("smtp connection refused")
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.AdminTicketEmailPayload{
			TicketID: ticketID})
		task := asynq.NewTask(tasks.TaskTypeAdminTicketEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.EqualError(t, err, "smtp connection refused")
	})

	t.Run("Success - Finished Ticket Email", func(t *testing.T) {
		ticketID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendFinishedTicketEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				assert.Equal(t, ticketID, receivedID)
				return nil
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.FinishedTicketEmailPayload{
			TicketID: ticketID})
		task := asynq.NewTask(tasks.TaskTypeFinishedTicketEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)
		assert.NoError(t, err)
	})

	t.Run("Error - Finished Ticket Email JSON Unmarshal", func(t *testing.T) {
		mockService := &mocks.MockEmailService{}
		handler := workers.NewEmailTaskHandler(mockService)

		task := asynq.NewTask(tasks.TaskTypeFinishedTicketEmail,
			[]byte(`[1, 2, 3]`))

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.ErrorIs(t, err, asynq.SkipRetry)
		assert.ErrorContains(t, err, "json unmarshal failed:")
	})

	t.Run("Error - Finished Ticket Email Service Failure", func(t *testing.T) {
		ticketID := uuid.New()
		mockService := &mocks.MockEmailService{
			SendFinishedTicketEmailFunc: func(ctx context.Context,
				receivedID uuid.UUID) error {
				return errors.New("template render error")
			},
		}
		handler := workers.NewEmailTaskHandler(mockService)

		payloadBytes, _ := json.Marshal(tasks.FinishedTicketEmailPayload{
			TicketID: ticketID})
		task := asynq.NewTask(tasks.TaskTypeFinishedTicketEmail, payloadBytes)

		err := handler.ProcessTask(ctx, task)

		assert.Error(t, err)
		assert.EqualError(t, err, "template render error")
	})
}
