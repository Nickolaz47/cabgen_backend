package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetEvents(ctx context.Context, limit int) ([]models.Event, error)
	CreateEvent(ctx context.Context, event *models.Event) error
	MarkProcessing(ctx context.Context, ID uint) error
	MarkDone(ctx context.Context, ID uint) error
	MarkFailed(ctx context.Context, ID uint, errorMessage string) error
}

type eventRepo struct {
	DB *gorm.DB
}

func NewEventRepo(db *gorm.DB) EventRepository {
	return &eventRepo{DB: db}
}

func (r *eventRepo) GetEvents(ctx context.Context, limit int) ([]models.Event, error) {
	var events []models.Event

	tx := r.DB.WithContext(ctx).Begin()

	if err := tx.Where("status = ?", models.EventPending).
		Order("created_at asc").
		Limit(limit).
		Find(&events).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(events) == 0 {
		return events, nil
	}

	if err := tx.Model(&events).
		Update("status", models.EventProcessing).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepo) CreateEvent(ctx context.Context, event *models.Event) error {
	return r.DB.WithContext(ctx).Create(event).Error
}

func (r *eventRepo) MarkProcessing(ctx context.Context, ID uint) error {
	return r.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", ID).
		Update("status", models.EventProcessing).Error
}

func (r *eventRepo) MarkDone(ctx context.Context, ID uint) error {
	return r.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", ID).
		Update("status", models.EventDone).Error
}

func (r *eventRepo) MarkFailed(ctx context.Context, ID uint,
	errorMessage string) error {
	return r.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", ID).
		Updates(map[string]any{
			"status": models.EventFailed,
			"error":  errorMessage,
		}).Error
}
