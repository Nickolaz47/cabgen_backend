package mocks

import (
	"context"
	"errors"

	"github.com/google/uuid"
	gomail "gopkg.in/mail.v2"
)

type MockEmailSender struct {
	ShouldFail bool
}

func (m *MockEmailSender) Send(msg *gomail.Message) error {
	if m.ShouldFail {
		return errors.New("simulated send error")
	}
	return nil
}

type MockEmailService struct {
	SendAdminAlertEmailFunc     func(ctx context.Context, newUserID uuid.UUID) error
	SendWelcomeEmailFunc        func(ctx context.Context, userID uuid.UUID) error
	SendAnalysisDoneEmailFunc   func(ctx context.Context, analysisID uuid.UUID) error
	SendAdminTicketEmailFunc    func(ctx context.Context, ticketID uuid.UUID) error
	SendFinishedTicketEmailFunc func(ctx context.Context, ticketID uuid.UUID) error
}

func (m *MockEmailService) SendAdminAlertEmail(ctx context.Context,
	newUserID uuid.UUID) error {
	if m.SendAdminAlertEmailFunc != nil {
		return m.SendAdminAlertEmailFunc(ctx, newUserID)
	}
	return nil
}

func (m *MockEmailService) SendWelcomeEmail(ctx context.Context,
	userID uuid.UUID) error {
	if m.SendWelcomeEmailFunc != nil {
		return m.SendWelcomeEmailFunc(ctx, userID)
	}
	return nil
}

func (m *MockEmailService) SendAnalysisDoneEmail(ctx context.Context,
	analysisID uuid.UUID) error {
	if m.SendAnalysisDoneEmailFunc != nil {
		return m.SendAnalysisDoneEmailFunc(ctx, analysisID)
	}
	return nil
}

func (m *MockEmailService) SendAdminTicketEmail(ctx context.Context,
	ticketID uuid.UUID) error {
	if m.SendAdminTicketEmailFunc != nil {
		return m.SendAdminTicketEmailFunc(ctx, ticketID)
	}
	return nil
}

func (m *MockEmailService) SendFinishedTicketEmail(ctx context.Context,
	ticketID uuid.UUID) error {
	if m.SendFinishedTicketEmailFunc != nil {
		return m.SendFinishedTicketEmailFunc(ctx, ticketID)
	}
	return nil
}
