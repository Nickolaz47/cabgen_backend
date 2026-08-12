package services

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/email"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
)

type EmailService interface {
	SendAdminAlertEmail(ctx context.Context, newUserID uuid.UUID) error
	SendWelcomeEmail(ctx context.Context, userID uuid.UUID) error
	SendAnalysisDoneEmail(ctx context.Context, analysisID uuid.UUID) error
	SendAdminTicketEmail(ctx context.Context, ticketID uuid.UUID) error
	SendFinishedTicketEmail(ctx context.Context, ticketID uuid.UUID) error
	SendPasswordResetEmail(ctx context.Context, userEmail, userName,
		token string) error
	SendUserDeletedEmail(ctx context.Context, userEmail, userName string) error
	SendEmailUpdateConfirmation(ctx context.Context, userEmail, userName,
		oldEmail, newEmail, token string) error
}

type emailService struct {
	UserRepo     repositories.UserRepository
	AnalysisRepo repositories.AnalysisRepository
	TicketRepo   repositories.TicketRepository
	EmailSender  email.EmailSender
	Logger       *zap.Logger
}

func NewEmailService(
	userRepo repositories.UserRepository,
	analysisRepo repositories.AnalysisRepository,
	ticketRepo repositories.TicketRepository,
	emailSender email.EmailSender,
	logger *zap.Logger) EmailService {
	return &emailService{
		UserRepo:     userRepo,
		AnalysisRepo: analysisRepo,
		TicketRepo:   ticketRepo,
		EmailSender:  emailSender,
		Logger:       logger,
	}
}

func (s *emailService) getLocalizer(language string) *i18n.Localizer {
	return i18n.NewLocalizer(translation.Bundle, language)
}

func (s *emailService) localize(localizer *i18n.Localizer, messageID string,
	data map[string]any) string {
	cfg := &i18n.LocalizeConfig{MessageID: messageID}
	if data != nil {
		cfg.TemplateData = data
	}

	msg, err := localizer.Localize(cfg)
	if err != nil {
		s.Logger.Error("Missing translation key",
			zap.String("messageID", messageID),
			zap.Error(err))
		return messageID
	}

	return msg
}

func (s *emailService) SendAdminAlertEmail(ctx context.Context,
	newUserID uuid.UUID) error {
	newUser, err := s.UserRepo.GetUserByID(ctx, newUserID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendAdminAlertEmail", logging.DatabaseError, err,
		)...)
		return fmt.Errorf("Failed to fetch new user: %v", err)
	}

	admin, isActive := models.Admin, true
	filter := models.AdminUserFilter{
		UserRole: admin,
		Active:   &isActive,
	}

	admins, err := s.UserRepo.GetUsers(ctx, filter)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendActivationUserEmail",
			logging.DatabaseError, err,
		)...)
		return fmt.Errorf("Failed to get admins: %v", err)
	}

	for _, a := range admins {
		if a.Email == "" {
			continue
		}

		localizer := s.getLocalizer(a.Language)
		subject := s.localize(localizer, "email.admin_alert.subject",
			map[string]any{"Username": newUser.Username})
		body := s.localize(localizer, "email.admin_alert.body",
			map[string]any{"Username": newUser.Username})

		cfg := email.EmailConfig{
			Sender:    config.SenderEmail,
			Recipient: a.Email,
			Subject:   subject,
			Body:      body,
		}
		if err := email.SendEmail(cfg, s.EmailSender); err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"EmailService", "SendAdminAlertEmail", logging.SendEmailError,
				fmt.Errorf("Failed to send alert to %s: %v", a.Email, err),
			)...)
		} else {
			s.Logger.Info("Email sent", logging.ServiceInfoLogging(
				"EmailService", "SendAdminAlertEmail", logging.EmailSentSuccess,
				zap.String("recipient", a.Email),
			)...)
		}
	}

	return nil
}

func (s *emailService) SendWelcomeEmail(ctx context.Context,
	userID uuid.UUID) error {
	user, err := s.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendWelcomeEmail", logging.DatabaseError, err,
		)...)
		return fmt.Errorf("Failed to fetch user: %v", err)
	}

	localizer := s.getLocalizer(user.Language)
	subject := s.localize(localizer, "email.welcome.subject", nil)
	body := s.localize(localizer, "email.welcome.body", map[string]any{
		"Name": user.Name,
	})

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: user.Email,
		Subject:   subject,
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendWelcomeEmail", logging.SendEmailError,
			fmt.Errorf("Failed to send welcome to %s: %v", user.Email, err),
		)...)
		return fmt.Errorf("Failed to send welcome email to %s: %v", user.Email,
			err)
	}

	s.Logger.Info("Email sent", logging.ServiceInfoLogging(
		"EmailService", "SendWelcomeEmail", logging.EmailSentSuccess,
		zap.String("recipient", user.Email),
	)...)

	return nil
}

func (s *emailService) SendAnalysisDoneEmail(ctx context.Context,
	analysisID uuid.UUID) error {
	analysis, err := s.AnalysisRepo.GetAnalysisByID(ctx, analysisID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendAnalysisDoneEmail", logging.DatabaseError, err,
		)...)
		return fmt.Errorf("Failed to fetch analysis: %v", err)
	}

	localizer := s.getLocalizer(analysis.User.Language)
	subject := s.localize(localizer, "email.analysis_done.subject", nil)

	statusText := s.localize(localizer, "email.analysis_done.status_done", nil)
	if analysis.Status == models.AnalysisStatusFailed {
		statusText = s.localize(localizer, "email.analysis_done.status_failed", nil)
	}

	body := s.localize(localizer, "email.analysis_done.body", map[string]any{
		"Name":             analysis.User.Name,
		"SampleOriginCode": analysis.Sample.OriginCode,
		"StatusText":       statusText,
	})

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: analysis.User.Email,
		Subject:   subject,
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendAnalysisDoneEmail", logging.SendEmailError,
			fmt.Errorf("Failed to send analysis email to %s: %v",
				analysis.User.Email, err),
		)...)
		return fmt.Errorf("Failed to send analysis email to %s: %v",
			analysis.User.Email, err)
	}

	s.Logger.Info("Email sent", logging.ServiceInfoLogging(
		"EmailService", "SendAnalysisDoneEmail", logging.EmailSentSuccess,
		zap.String("recipient", analysis.User.Email),
		zap.String("analysis_id", analysisID.String()),
	)...)

	return nil
}

func (s *emailService) SendAdminTicketEmail(ctx context.Context,
	ticketID uuid.UUID) error {
	ticket, err := s.TicketRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendAdminTicketEmail", logging.DatabaseError, err,
		)...)
		return fmt.Errorf("Failed to fetch ticket: %v", err)
	}

	admin, isActive := models.Admin, true
	filter := models.AdminUserFilter{
		UserRole: admin,
		Active:   &isActive,
	}

	admins, err := s.UserRepo.GetUsers(ctx, filter)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendAdminTicketEmail",
			logging.DatabaseError, err,
		)...)
		return fmt.Errorf("Failed to get admins: %v", err)
	}

	for _, a := range admins {
		if a.Email == "" {
			continue
		}

		localizer := s.getLocalizer(a.Language)
		subject := s.localize(localizer, "email.admin_ticket.subject",
			map[string]any{"Name": ticket.Name})
		body := s.localize(localizer, "email.admin_ticket.body",
			map[string]any{
				"Name":    ticket.Name,
				"Email":   ticket.Email,
				"Subject": ticket.Subject,
				"Message": ticket.Message,
			})

		cfg := email.EmailConfig{
			Sender:    config.SenderEmail,
			Recipient: a.Email,
			Subject:   subject,
			Body:      body,
		}

		if err := email.SendEmail(cfg, s.EmailSender); err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"EmailService", "SendAdminTicketEmail", logging.SendEmailError,
				fmt.Errorf("Failed to send ticket email to %s: %v",
					a.Email, err),
			)...)
		} else {
			s.Logger.Info("Email sent", logging.ServiceInfoLogging(
				"EmailService", "SendAdminTicketEmail", logging.EmailSentSuccess,
				zap.String("recipient", a.Email),
			)...)
		}
	}

	return nil
}

func (s *emailService) SendFinishedTicketEmail(ctx context.Context,
	ticketID uuid.UUID) error {
	ticket, err := s.TicketRepo.GetTicketByID(ctx, ticketID)
	if err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendFinishedTicketEmail", logging.DatabaseError,
			err)...)
		return fmt.Errorf("Failed to fetch ticket: %v", err)
	}

	localizer := s.getLocalizer(ticket.Language)
	subject := s.localize(localizer, "email.finished_ticket.subject",
		map[string]any{"Subject": ticket.Subject})
	body := s.localize(localizer, "email.finished_ticket.body",
		map[string]any{
			"Name":    ticket.Name,
			"Subject": ticket.Subject,
			"Message": ticket.Message,
		})

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: ticket.Email,
		Subject:   subject,
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendFinishedTicketEmail", logging.SendEmailError,
			fmt.Errorf("Failed to send ticket email to %s: %v", ticket.Email,
				err),
		)...)
		return fmt.Errorf("Failed to send ticket email to %s: %v",
			ticket.Email, err)
	}

	s.Logger.Info("Email sent", logging.ServiceInfoLogging(
		"EmailService", "SendFinishedTicketEmail", logging.EmailSentSuccess,
		zap.String("recipient", ticket.Email),
		zap.String("ticket_id", ticketID.String()),
	)...)

	return nil
}

func (s *emailService) SendPasswordResetEmail(ctx context.Context, userEmail,
	userName, token string) error {
	frontendURL := config.FrontendURL
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	language := "en"
	if s.UserRepo != nil {
		if user, err := s.UserRepo.GetUserByEmail(ctx, userEmail); err == nil {
			language = user.Language
		}
	}

	localizer := s.getLocalizer(language)
	subject := s.localize(localizer, "email.password_reset.subject", nil)
	body := s.localize(localizer, "email.password_reset.body", map[string]any{
		"Name":      userName,
		"ResetLink": resetLink,
	})

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: userEmail,
		Subject:   subject,
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendPasswordResetEmail", logging.SendEmailError,
			fmt.Errorf("Failed to send reset email to %s: %v", userEmail, err),
		)...)
		return fmt.Errorf("Failed to send reset email to %s: %v", userEmail, err)
	}

	s.Logger.Info("Email sent", logging.ServiceInfoLogging(
		"EmailService", "SendPasswordResetEmail", logging.EmailSentSuccess,
		zap.String("recipient", userEmail),
	)...)

	return nil
}

func (s *emailService) SendUserDeletedEmail(ctx context.Context, userEmail,
	userName string) error {
	language := "en"
	if s.UserRepo != nil {
		if user, err := s.UserRepo.GetUserByEmail(ctx, userEmail); err == nil {
			language = user.Language
		}
	}

	localizer := s.getLocalizer(language)
	subject := s.localize(localizer, "email.user_deleted.subject", nil)
	body := s.localize(localizer, "email.user_deleted.body", map[string]any{
		"Name": userName,
	})

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: userEmail,
		Subject:   subject,
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendUserDeletedEmail", logging.SendEmailError,
			fmt.Errorf("Failed to send deletion email to %s: %v", userEmail, err),
		)...)
		return fmt.Errorf("Failed to send deletion email to %s: %v", userEmail, err)
	}

	s.Logger.Info("Email sent", logging.ServiceInfoLogging(
		"EmailService", "SendUserDeletedEmail", logging.EmailSentSuccess,
		zap.String("recipient", userEmail),
	)...)

	return nil
}

func (s *emailService) SendEmailUpdateConfirmation(ctx context.Context,
	userEmail, userName, oldEmail, newEmail, token string) error {
	frontendURL := config.FrontendURL
	confirmLink := fmt.Sprintf("%s/confirm-email-update?token=%s",
		frontendURL, token)

	language := "en"
	if s.UserRepo != nil {
		if user, err := s.UserRepo.GetUserByEmail(ctx, userEmail); err == nil {
			language = user.Language
		}
	}

	localizer := s.getLocalizer(language)
	subject := s.localize(localizer, "email.email_update_confirmation.subject", nil)
	body := s.localize(localizer, "email.email_update_confirmation.body",
		map[string]any{
			"Name":        userName,
			"OldEmail":    oldEmail,
			"NewEmail":    newEmail,
			"ConfirmLink": confirmLink,
		})

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: newEmail,
		Subject:   subject,
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendEmailUpdateConfirmation",
			logging.SendEmailError,
			fmt.Errorf("Failed to send email update confirmation to %s: %v",
				newEmail, err),
		)...)
		return fmt.Errorf("Failed to send email update confirmation to %s: %v",
			newEmail, err)
	}

	s.Logger.Info("Email sent", logging.ServiceInfoLogging(
		"EmailService", "SendEmailUpdateConfirmation", logging.EmailSentSuccess,
		zap.String("recipient", newEmail),
		zap.String("old_email", oldEmail),
	)...)

	return nil
}
