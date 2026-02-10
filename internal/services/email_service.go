package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/email"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"go.uber.org/zap"
)

type EmailService interface {
	SendActivationUserEmail(ctx context.Context, userToActivate string) error
}

type emailService struct {
	UserRepo    repositories.UserRepository
	EmailSender email.EmailSender
	Logger      *zap.Logger
}

func NewEmailService(
	userRepo repositories.UserRepository,
	emailSender email.EmailSender,
	logger *zap.Logger) EmailService {
	return &emailService{
		UserRepo:    userRepo,
		EmailSender: emailSender,
		Logger:      logger,
	}
}

func (s *emailService) SendActivationUserEmail(
	ctx context.Context, userToActivate string) error {
	admin, isActive := models.Admin, true
	filter := models.AdminUserFilter{
		UserRole: &admin,
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

	var adminEmailConfigs []email.EmailConfig
	for _, a := range admins {
		body := `
		Prezados administradores,
		<p>Um novo usuário foi criado no CAGBen.</p>
		<p>Por favor, acesse o site e realize a ativação do mesmo.</p>
		Obrigado.
		`
		if a.Email != "" {
			emailConfig := email.EmailConfig{
				Sender:    config.SenderEmail,
				Recipient: a.Email,
				Subject:   "Novo Usuário Criado: " + userToActivate,
				Body:      body,
			}
			adminEmailConfigs = append(adminEmailConfigs, emailConfig)
		}
	}

	var wg sync.WaitGroup

	for _, emailConfig := range adminEmailConfigs {
		wg.Add(1)

		go func(cfg email.EmailConfig) {
			defer wg.Done()

			err := email.SendEmail(cfg, s.EmailSender)
			if err != nil {
				s.Logger.Error(
					"Service Error", logging.ServiceLogging(
						"EmailService", "SendActivationUserEmail",
						logging.SendEmailError,
						fmt.Errorf(
							"Failed to send activation user email to %s: %v",
							cfg.Recipient, err),
					)...)
			}
		}(emailConfig)
	}
	wg.Wait()

	return nil
}
