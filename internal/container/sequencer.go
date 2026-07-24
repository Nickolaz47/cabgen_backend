package container

import (
	adminSequencer "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BuildSequencerService(db *gorm.DB, logger *zap.Logger) services.SequencerService {
	sequencerRepo := repositories.NewSequencerRepo(db)
	sequencerService := services.NewSequencerService(sequencerRepo, logger)

	return sequencerService
}

func BuildAdminSequencerHandler(svc services.SequencerService) *adminSequencer.AdminSequencerHandler {
	return adminSequencer.NewAdminSequencerHandler(svc)
}
