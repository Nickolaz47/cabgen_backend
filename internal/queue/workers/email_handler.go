package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
)

type EmailTaskHandler struct {
	EmailService services.EmailService
}

func NewEmailTaskHandler(emailService services.EmailService) *EmailTaskHandler {
	return &EmailTaskHandler{
		EmailService: emailService,
	}
}

func (h *EmailTaskHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	switch t.Type() {

	case tasks.TaskTypeAdminAlertEmail:
		var p tasks.AdminAlertEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		return h.EmailService.SendAdminAlertEmail(ctx, p.NewUserID)

	case tasks.TaskTypeWelcomeEmail:
		var p tasks.WelcomeEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		return h.EmailService.SendWelcomeEmail(ctx, p.UserID)

	case tasks.TaskTypeAnalysisDoneEmail:
		var p tasks.AnalysisDoneEmailPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}
		return h.EmailService.SendAnalysisDoneEmail(ctx, p.AnalysisID)

	default:
		return fmt.Errorf("unknown task type: %s", t.Type())
	}
}
