package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type EmailTaskHandler struct {
	EmailService services.EmailService
	Logger       *zap.Logger
}

func NewEmailTaskHandler(emailService services.EmailService,
	logger *zap.Logger) *EmailTaskHandler {
	return &EmailTaskHandler{
		EmailService: emailService,
		Logger:       logger,
	}
}

func (h *EmailTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	switch t.Type() {

	case tasks.TaskTypeAdminAlertEmail:
		var p tasks.AdminAlertEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("new_user_id", p.NewUserID.String()))
		return h.execute(t, h.EmailService.SendAdminAlertEmail(ctx, p.NewUserID))

	case tasks.TaskTypeWelcomeEmail:
		var p tasks.WelcomeEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("user_id", p.UserID.String()))
		return h.execute(t, h.EmailService.SendWelcomeEmail(ctx, p.UserID))

	case tasks.TaskTypeAnalysisDoneEmail:
		var p tasks.AnalysisDoneEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("analysis_id", p.AnalysisID.String()))
		return h.execute(t, h.EmailService.SendAnalysisDoneEmail(ctx, p.AnalysisID))

	case tasks.TaskTypeAdminTicketEmail:
		var p tasks.AdminTicketEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("ticket_id", p.TicketID.String()))
		return h.execute(t, h.EmailService.SendAdminTicketEmail(ctx, p.TicketID))

	case tasks.TaskTypeFinishedTicketEmail:
		var p tasks.FinishedTicketEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("ticket_id", p.TicketID.String()))
		return h.execute(t, h.EmailService.SendFinishedTicketEmail(ctx, p.TicketID))

	case tasks.TaskTypePasswordResetEmail:
		var p tasks.PasswordResetEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("email", p.Email))
		return h.execute(t, h.EmailService.SendPasswordResetEmail(ctx, p.Email,
			p.Name, p.Token))

	case tasks.TaskTypeUserDeletedEmail:
		var p tasks.UserDeletedEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("email", p.Email))
		return h.execute(t, h.EmailService.SendUserDeletedEmail(ctx, p.Email, p.Name))

	case tasks.TaskTypeEmailUpdateConfirmation:
		var p tasks.EmailUpdateConfirmationPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		h.logTaskStart(t, zap.String("email", p.Email),
			zap.String("new_email", p.NewEmail))
		return h.execute(t, h.EmailService.SendEmailUpdateConfirmation(ctx, p.Email,
			p.Name, p.OldEmail, p.NewEmail, p.Token))

	default:
		return fmt.Errorf("unknown task type: %s", t.Type())
	}
}

func (h *EmailTaskHandler) logTaskStart(t *asynq.Task, fields ...zap.Field) {
	h.Logger.Info("Task started", logging.ServiceInfoLogging(
		"EmailTaskHandler", "ProcessTask", "TASK_STARTED",
		append([]zap.Field{zap.String("task_type", t.Type())}, fields...)...,
	)...)
}

func (h *EmailTaskHandler) execute(t *asynq.Task, err error) error {
	if err != nil {
		h.Logger.Error("Task failed", logging.ServiceLogging(
			"EmailTaskHandler", "ProcessTask", logging.SendEmailError, err)...)
		return err
	}

	h.Logger.Info("Task completed", logging.ServiceInfoLogging(
		"EmailTaskHandler", "ProcessTask", "TASK_COMPLETED",
		zap.String("task_type", t.Type()),
	)...)
	return nil
}
