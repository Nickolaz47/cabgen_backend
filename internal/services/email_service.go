package services

import (
	"context"
	"fmt"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/email"
	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EmailService interface {
	SendAdminAlertEmail(ctx context.Context, newUserID uuid.UUID) error
	SendWelcomeEmail(ctx context.Context, userID uuid.UUID) error
	SendAnalysisDoneEmail(ctx context.Context, analysisID uuid.UUID) error
}

type emailService struct {
	UserRepo     repositories.UserRepository
	AnalysisRepo repositories.AnalysisRepository
	EmailSender  email.EmailSender
	Logger       *zap.Logger
}

func NewEmailService(
	userRepo repositories.UserRepository,
	analysisRepo repositories.AnalysisRepository,
	emailSender email.EmailSender,
	logger *zap.Logger) EmailService {
	return &emailService{
		UserRepo:     userRepo,
		AnalysisRepo: analysisRepo,
		EmailSender:  emailSender,
		Logger:       logger,
	}
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

	for _, a := range admins {
		if a.Email == "" {
			continue
		}

		body := `
		Prezado administrador,
		<p>Um novo usuário foi criado no CAGBen.</p>
		<p>Por favor, acesse o site e realize a ativação do mesmo.</p>
		Obrigado.
		`

		cfg := email.EmailConfig{
			Sender:    config.SenderEmail,
			Recipient: a.Email,
			Subject:   "Novo Usuário Criado: " + newUser.Username,
			Body:      body,
		}
		if err := email.SendEmail(cfg, s.EmailSender); err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"EmailService", "SendAdminAlertEmail", logging.SendEmailError,
				fmt.Errorf("Failed to send alert to %s: %v", a.Email, err),
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

	body := fmt.Sprintf(`
	Olá, %s!
	<p>Sua conta no CABGen acaba de ser ativada por um administrador.</p>
	<p>Você já pode realizar o login e começar a analisar suas amostras.</p>
	<br>Equipe CABGen.
	`, user.Name)

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: user.Email,
		Subject:   "Sua conta CABGen foi ativada!",
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

	statusText := "foi finalizada com sucesso"
	if analysis.Status == models.AnalysisStatusFailed {
		statusText = "encontrou um erro durante o processamento"
	}

	body := fmt.Sprintf(`
	Olá, %s!
	<p>A análise da sua amostra <strong>%s</strong> %s.</p>
	<p>Acesse o sistema para verificar os resultados detalhados.</p>
	<br>Equipe CABGen.
	`, analysis.User.Name, analysis.Sample.Name, statusText)

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: analysis.User.Email,
		Subject:   "CABGen - Análise Finalizada",
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

	return nil
}
