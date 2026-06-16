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

	TaskTypeAnalysisProcess     = "analysis:process"
	TaskTypeWelcomeEmail        = "email:welcome"
	TaskTypeAnalysisDoneEmail   = "email:analysis_done"
	TaskTypeAdminAlertEmail     = "email:admin_user_alert"
	TaskTypeAdminTicketEmail    = "email:admin_ticket"
	TaskTypeFinishedTicketEmail = "email:finished_ticket"
	TaskTypePasswordResetEmail  = "email:password_reset"
)

type AnalysisProcessPayload struct {
	AnalysisID uuid.UUID `json:"analysis_id"`
}

type WelcomeEmailPayload struct {
	UserID uuid.UUID `json:"user_id"`
}

type AnalysisDoneEmailPayload struct {
	AnalysisID uuid.UUID `json:"analysis_id"`
}

type AdminAlertEmailPayload struct {
	NewUserID uuid.UUID `json:"new_user_id"`
}

type AdminTicketEmailPayload struct {
	TicketID uuid.UUID `json:"ticket_id"`
}

type FinishedTicketEmailPayload struct {
	TicketID uuid.UUID `json:"ticket_id"`
}

type PasswordResetEmailPayload struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

func NewAnalysisProcessTask(analysisID uuid.UUID) (
	*asynq.Task, error) {
	payload := AnalysisProcessPayload{AnalysisID: analysisID}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskTypeAnalysisProcess,
		payloadBytes,
		asynq.MaxRetry(3),
		asynq.Timeout(5*time.Hour),
	), nil
}

func NewAdminAlertEmailTask(newUserID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(AdminAlertEmailPayload{NewUserID: newUserID})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskTypeAdminAlertEmail, payload), nil
}

func NewWelcomeEmailTask(userID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{UserID: userID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeWelcomeEmail, payload,
		asynq.MaxRetry(5)), nil
}

func NewAnalysisDoneEmailTask(analysisID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(AnalysisDoneEmailPayload{
		AnalysisID: analysisID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeAnalysisDoneEmail, payload,
		asynq.MaxRetry(5)), nil
}

func NewAdminTicketEmailTask(ticketID uuid.UUID) (
	*asynq.Task, error) {
	payload, err := json.Marshal(AdminTicketEmailPayload{
		TicketID: ticketID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeAdminTicketEmail, payload, asynq.MaxRetry(5)),
		nil
}

func NewFinishedTicketEmailTask(ticketID uuid.UUID) (*asynq.Task, error) {
	payload, err := json.Marshal(FinishedTicketEmailPayload{
		TicketID: ticketID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypeFinishedTicketEmail, payload,
		asynq.MaxRetry(5)), nil
}

func NewPasswordResetEmailTask(email, name, token string) (*asynq.Task, error) {
	payload, err := json.Marshal(PasswordResetEmailPayload{
		Email: email,
		Name:  name,
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskTypePasswordResetEmail, payload,
		asynq.MaxRetry(3)), nil
}
