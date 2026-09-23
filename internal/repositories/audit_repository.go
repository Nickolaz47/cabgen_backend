package repositories

import (
	"context"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

type AuditRepository interface {
	GetAuditLogs(ctx context.Context, auditFilter models.AuditFilter) (
		[]models.Audit, error)
	CreateAudit(ctx context.Context, audit *models.Audit) error
}

type auditRepo struct {
	DB *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepo{DB: db}
}

func (r *auditRepo) GetAuditLogs(ctx context.Context,
	auditFilter models.AuditFilter) ([]models.Audit, error) {
	var auditLogs []models.Audit

	query := r.DB.WithContext(ctx).Preload("User")
	if auditFilter.Event != "" {
		query = query.Where("event = ?", auditFilter.Event)
	}

	if auditFilter.Source != "" {
		like := "%" + auditFilter.Source + "%"
		query = query.Where("source LIKE ?", like)
	}

	if auditFilter.Status != 0 {
		query = query.Where("status = ?", auditFilter.Status)
	}

	if auditFilter.Date != nil {
		start := *auditFilter.Date
		end := start.Add(24 * time.Hour)
		query = query.Where("created_at >= ? AND created_at < ?", start, end)
	}

	if auditFilter.UserID != nil {
		query = query.Where("user_id = ?", auditFilter.UserID)
	}

	if err := query.Order("created_at DESC").
		Find(&auditLogs).Error; err != nil {
		return nil, err
	}

	return auditLogs, nil
}

func (r *auditRepo) CreateAudit(ctx context.Context, audit *models.Audit) error {
	return r.DB.WithContext(ctx).Create(audit).Error
}
