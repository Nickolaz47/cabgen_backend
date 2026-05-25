package repositories

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetEvents(ctx context.Context, limit int) ([]models.Event, error)
	CreateEvent(ctx context.Context, event *models.Event) error
	MarkProcessing(ctx context.Context, ID uuid.UUID) error
	MarkDone(ctx context.Context, ID uuid.UUID) error
	MarkFailed(ctx context.Context, ID uuid.UUID, errorMessage string) error
}

type eventRepo struct {
	DB *gorm.DB
}

func NewEventRepo(db *gorm.DB) EventRepository {
	return &eventRepo{DB: db}
}

func (r *eventRepo) GetEvents(ctx context.Context,
	limit int) ([]models.Event, error) {
	var events []models.Event

	dialector := r.DB.Dialector.Name()

	if dialector == "postgres" {
		query := `
            UPDATE events SET status = ?
            WHERE id IN (
                SELECT id FROM events
                WHERE status = ?
                ORDER BY created_at ASC
                LIMIT ?
                FOR UPDATE SKIP LOCKED
            )
            RETURNING *
        `
		err := r.DB.WithContext(ctx).
			Raw(query, models.EventProcessing, models.EventPending, limit).
			Scan(&events).Error
		if err != nil {
			return nil, err
		}
		return events, nil
	}

	// Test only
	err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Event{}).
			Where("status = ?", models.EventPending).
			Order("created_at ASC").
			Limit(limit).
			Find(&events).Error; err != nil {
			return err
		}

		if len(events) == 0 {
			return nil
		}

		return tx.Model(&events).
			Update("status", models.EventProcessing).Error
	})

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepo) CreateEvent(ctx context.Context, event *models.Event) error {
	return r.DB.WithContext(ctx).Create(event).Error
}

func (r *eventRepo) MarkProcessing(ctx context.Context, ID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", ID).
		Update("status", models.EventProcessing).Error
}

func (r *eventRepo) MarkDone(ctx context.Context, ID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", ID).
		Update("status", models.EventDone).Error
}

func (r *eventRepo) MarkFailed(ctx context.Context, ID uuid.UUID,
	errorMessage string) error {
	return r.DB.WithContext(ctx).Model(&models.Event{}).Where("id = ?", ID).
		Updates(map[string]any{
			"status": models.EventFailed,
			"error":  errorMessage,
		}).Error
}
