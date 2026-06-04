package tasks

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	QueueAnalysis = "analyses"
	QueueEmail    = "emails"

	TaskTypeProcessAnalysis    = "analysis:process"
	TaskTypeRegisterUserEmail  = "email:register"
	TaskTypeUserActivatedEmail = "email:activation"
	TaskTypeAnalysisDoneEmail  = "email:analysisDone"
)

type AnalysisProcessPayload struct {
	AnalysisID uuid.UUID `json:"analysis_id"`
}

type UserEmailPayload struct {
	UserID uuid.UUID `json:"user_id"`
}

type AnalysisDoneEmailPayload struct {
	AnalysisID uuid.UUID `json:"analysis_id"`
}

func NewProcessAnalysisTask(analysisID uuid.UUID) (
	*asynq.Task, error) {
	payload := AnalysisProcessPayload{AnalysisID: analysisID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskTypeProcessAnalysis,
		payloadBytes,
		asynq.MaxRetry(3),
		asynq.Timeout(5*time.Hour),
	), nil
}

func NewRegisterUserEmailTask(userID uuid.UUID) (*asynq.Task, error) {
	payload := UserEmailPayload{UserID: userID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskTypeRegisterUserEmail,
		payloadBytes,
		asynq.MaxRetry(5),
	), nil
}

func NewUserActivatedEmailTask(userID uuid.UUID) (*asynq.Task, error) {
	payload := UserEmailPayload{UserID: userID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskTypeUserActivatedEmail,
		payloadBytes,
		asynq.MaxRetry(5),
	), nil
}

func NewAnalysisDoneEmailTask(analysisID uuid.UUID) (*asynq.Task, error) {
	payload := AnalysisDoneEmailPayload{AnalysisID: analysisID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskTypeAnalysisDoneEmail,
		payloadBytes,
		asynq.MaxRetry(5),
	), nil
}
