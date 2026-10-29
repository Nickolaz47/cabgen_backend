package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/email"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildEmailService(db *gorm.DB, logger *zap.Logger) services.EmailService {
	userRepo := repositories.NewUserRepo(db)
	analysisRepo := repositories.NewAnalysisRepository(db)
	ticketRepo := repositories.NewTicketRepo(db)
	emailSvc := services.NewEmailService(
		userRepo, analysisRepo, ticketRepo,
		repositories.NewPasswordResetRepo(db),
		repositories.NewEmailUpdateRepo(db),
		email.CreateDefaultSender(), logger,
	)

	return emailSvc
}
