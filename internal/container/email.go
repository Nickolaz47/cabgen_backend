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
	emailSvc := services.NewEmailService(
		userRepo, email.CreateDefaultSender(), logger,
	)

	return emailSvc
}
