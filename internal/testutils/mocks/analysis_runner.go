package mocks

import (
	"context"

	"github.com/google/uuid"
)

type MockAnalysisRunnerService struct {
	RunFunc func(ctx context.Context, analysisID uuid.UUID) error
}

func (s *MockAnalysisRunnerService) Run(ctx context.Context,
	analysisID uuid.UUID) error {
	if s.RunFunc != nil {
		return s.RunFunc(ctx, analysisID)
	}
	return nil
}