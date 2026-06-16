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
	SendAdminTicketEmail(ctx context.Context, ticketID uuid.UUID) error
	SendFinishedTicketEmail(ctx context.Context, ticketID uuid.UUID) error
	SendPasswordResetEmail(ctx context.Context, userEmail, userName,
		token string) error
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
		UserRole: &admin,
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

	body := fmt.Sprintf(`
    <h2>Novo Ticket de Suporte CABGen</h2>
    <p><strong>Nome:</strong> %s</p>
    <p><strong>E-mail:</strong> %s</p>
    <p><strong>Assunto:</strong> %s</p>
    <hr>
    <p>%s</p>
	<br>
	<p><small>Acesse o painel administrativo para atribuir este ticket a você e respondê-lo.</small></p>
    `, ticket.Name, ticket.Email, ticket.Subject, ticket.Message)

	for _, a := range admins {
		if a.Email == "" {
			continue
		}

		cfg := email.EmailConfig{
			Sender:    config.SenderEmail,
			Recipient: a.Email,
			Subject:   "Contato CABGen - " + ticket.Name,
			Body:      body,
		}

		if err := email.SendEmail(cfg, s.EmailSender); err != nil {
			s.Logger.Error("Service Error", logging.ServiceLogging(
				"EmailService", "SendAdminTicketEmail", logging.SendEmailError,
				fmt.Errorf("Failed to send ticket email to %s: %v",
					a.Email, err),
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

	body := fmt.Sprintf(`
    <h2>Atualização do seu Ticket de Suporte CABGen</h2>
    <p>Olá, <strong>%s</strong>,</p>
    <p>O seu chamado referente ao assunto "<strong>%s</strong>" foi analisado e marcado como <strong>Resolvido</strong> pela nossa equipe.</p>
    <hr>
    <p><small><strong>Lembrete da sua mensagem original:</strong><br>%s</small></p>
    <br>
    <p>Se o problema persistir ou se você tiver novas dúvidas, sinta-se à vontade para abrir um novo ticket em nosso site.</p>
    <p>Atenciosamente,<br><strong>Equipe CABGen</strong></p>
    `, ticket.Name, ticket.Subject, ticket.Message)

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: ticket.Email,
		Subject:   "CABGen - Ticket Resolvido: " + ticket.Subject,
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

	return nil
}

func (s *emailService) SendPasswordResetEmail(ctx context.Context, userEmail,
	userName, token string) error {
	frontendURL := config.FrontendURL
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	body := fmt.Sprintf(`
	<div style="font-family: Arial, sans-serif; color: #333;">
		<h2>Recuperação de Senha - CABGen</h2>
		<p>Olá, <strong>%s</strong>,</p>
		<p>Recebemos uma solicitação para redefinir a senha da sua conta no sistema CABGen.</p>
		<p>Se você não fez essa solicitação, por favor ignore este e-mail. Sua senha permanecerá a mesma.</p>
		<p>Para criar uma nova senha, clique no botão abaixo (este link expira em 15 minutos):</p>
		<br>
		<a href="%s" style="background-color: #0056b3; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; font-weight: bold;">Redefinir Minha Senha</a>
		<br><br>
		<p><small>Ou copie e cole o seguinte link no seu navegador:</small><br>
		<small><a href="%s">%s</a></small></p>
		<hr>
		<p>Atenciosamente,<br><strong>Equipe CABGen</strong></p>
	</div>
	`, userName, resetLink, resetLink, resetLink)

	cfg := email.EmailConfig{
		Sender:    config.SenderEmail,
		Recipient: userEmail,
		Subject:   "CABGen - Redefinição de Senha",
		Body:      body,
	}

	if err := email.SendEmail(cfg, s.EmailSender); err != nil {
		s.Logger.Error("Service Error", logging.ServiceLogging(
			"EmailService", "SendPasswordResetEmail", logging.SendEmailError,
			fmt.Errorf("Failed to send reset email to %s: %v", userEmail, err),
		)...)
		return fmt.Errorf("Failed to send reset email to %s: %v", userEmail, err)
	}

	return nil
}
