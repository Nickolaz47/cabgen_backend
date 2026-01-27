package mocks

import (
	"context"
	"errors"

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
	SendActivationUserEmailFunc func(
		ctx context.Context, userToActivate string) error
}

func (s *MockEmailService) SendActivationUserEmail(ctx context.Context,
	userToActivate string) error {
	if s.SendActivationUserEmailFunc != nil {
		return s.SendActivationUserEmailFunc(ctx, userToActivate)
	}
	return nil
}
