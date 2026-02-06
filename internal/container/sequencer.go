package container

import (
	adminSequencer "github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sequencer"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"gorm.io/gorm"
)

func BuildSequencerService(db *gorm.DB) services.SequencerService {
	sequencerRepo := repositories.NewSequencerRepo(db)
	sequencerService := services.NewSequencerService(sequencerRepo)

	return sequencerService
}

func BuildSequencerHandler(svc services.SequencerService) *sequencer.SequencerHandler {
	return sequencer.NewSequencerHandler(svc)
}

func BuildAdminSequencerHandler(svc services.SequencerService) *adminSequencer.AdminSequencerHandler {
	return adminSequencer.NewAdminSequencerHandler(svc)
}
