package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailUpdateRepository interface {
	CreateRequest(ctx context.Context, req *models.EmailUpdateRequest) error
	GetByToken(ctx context.Context, token string) (*models.EmailUpdateRequest, error)
	DeleteRequestsByUserID(ctx context.Context, userID uuid.UUID) error
}

type emailUpdateRepository struct {
	db *gorm.DB
}

func NewEmailUpdateRepo(db *gorm.DB) EmailUpdateRepository {
	return &emailUpdateRepository{db: db}
}

func (r *emailUpdateRepository) CreateRequest(ctx context.Context,
	req *models.EmailUpdateRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *emailUpdateRepository) GetByToken(ctx context.Context,
	token string) (*models.EmailUpdateRequest, error) {
	var req models.EmailUpdateRequest
	err := r.db.WithContext(ctx).First(&req, "token = ?", token).Error
	return &req, err
}

func (r *emailUpdateRepository) DeleteRequestsByUserID(ctx context.Context,
	userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(
		&models.EmailUpdateRequest{}).Error
}
