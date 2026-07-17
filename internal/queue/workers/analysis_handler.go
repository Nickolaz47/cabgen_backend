package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/queue/tasks"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/hibiken/asynq"
)

type AnalysisTaskHandler struct {
	AnalysisRunnerService services.AnalysisRunnerService
}

func NewAnalysisTaskHandler(
	analysisRunnerService services.AnalysisRunnerService) *AnalysisTaskHandler {
	return &AnalysisTaskHandler{
		AnalysisRunnerService: analysisRunnerService,
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

		return h.AnalysisRunnerService.Run(ctx, p.AnalysisID)
	default:
		return fmt.Errorf("unknown task type: %s", t.Type())
	}
}
