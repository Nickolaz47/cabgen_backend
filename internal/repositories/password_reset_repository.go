package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

type PasswordResetRepository interface {
	CreateToken(ctx context.Context, reset *models.PasswordReset) error
	GetByToken(ctx context.Context, token string) (*models.PasswordReset, error)
	DeleteTokensByEmail(ctx context.Context, email string) error
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepo(db *gorm.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) CreateToken(ctx context.Context,
	reset *models.PasswordReset) error {
	return r.db.WithContext(ctx).Create(reset).Error
}

func (r *passwordResetRepository) GetByToken(ctx context.Context,
	token string) (*models.PasswordReset, error) {
	var reset models.PasswordReset
	err := r.db.WithContext(ctx).First(&reset, "token = ?", token).Error
	return &reset, err
}

func (r *passwordResetRepository) DeleteTokensByEmail(ctx context.Context,
	email string) error {
	return r.db.WithContext(ctx).Where("email = ?", email).Delete(
		&models.PasswordReset{}).Error
}
