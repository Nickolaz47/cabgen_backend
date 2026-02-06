package container

import (
	"github.com/CABGenOrg/cabgen_backend/internal/email"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildEmailService(db *gorm.DB) services.EmailService {
	userRepo := repositories.NewUserRepo(db)
	emailSvc := services.NewEmailService(
		userRepo, email.CreateDefaultSender(),
	)

	return emailSvc
}
