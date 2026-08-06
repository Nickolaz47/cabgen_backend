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

type AnalysisTaskHandler struct {
	AnalysisRunnerService services.AnalysisRunnerService
	Logger                *zap.Logger
}

func NewAnalysisTaskHandler(
	analysisRunnerService services.AnalysisRunnerService,
	logger *zap.Logger) *AnalysisTaskHandler {
	return &AnalysisTaskHandler{
		AnalysisRunnerService: analysisRunnerService,
		Logger:                logger,
	}
}

func (h *AnalysisTaskHandler) ProcessTask(ctx context.Context,
	t *asynq.Task) error {
	switch t.Type() {
	case tasks.TaskTypeAnalysisProcess:
		var p tasks.AnalysisProcessPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", asynq.SkipRetry)
		}

		h.Logger.Info("Task started", logging.ServiceInfoLogging(
			"AnalysisTaskHandler", "ProcessTask", "TASK_STARTED",
			zap.String("task_type", t.Type()),
			zap.String("analysis_id", p.AnalysisID.String()),
		)...)

		if err := h.AnalysisRunnerService.Run(ctx, p.AnalysisID); err != nil {
			h.Logger.Error("Task failed", logging.ServiceLogging(
				"AnalysisTaskHandler", "ProcessTask", logging.AnalysisRunError,
				err)...)
			return err
		}

		h.Logger.Info("Task completed", logging.ServiceInfoLogging(
			"AnalysisTaskHandler", "ProcessTask", "TASK_COMPLETED",
			zap.String("task_type", t.Type()),
			zap.String("analysis_id", p.AnalysisID.String()),
		)...)
		return nil
	default:
		return fmt.Errorf("unknown task type: %s", t.Type())
	}
}
