package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewEventRepo(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)

	assert.NotEmpty(t, eventRepo)
}

func TestGetEvents(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	mockEvent := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "john"}`),
		Status:  models.EventPending,
	}
	mockEvent2 := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "maria"}`),
		Status:  models.EventDone,
	}
	mockEvent3 := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "ana"}`),
		Status:  models.EventPending,
	}

	db.Create(&mockEvent)
	db.Create(&mockEvent2)
	db.Create(&mockEvent3)

	t.Run("Success", func(t *testing.T) {
		result, err := eventRepo.GetEvents(ctx, 10)
		assert.NoError(t, err)

		var event1, event3 models.Event
		err = db.Where("id = ?", 1).First(&event1).Error
		assert.NoError(t, err)

		err = db.Where("id = ?", 3).First(&event3).Error
		assert.NoError(t, err)

		expectedEvent1 := models.Event{
			ID:        event1.ID,
			Name:      mockEvent.Name,
			Payload:   mockEvent.Payload,
			Status:    models.EventProcessing,
			CreatedAt: event1.CreatedAt,
		}
		expectedEvent3 := models.Event{
			ID:        event3.ID,
			Name:      mockEvent3.Name,
			Payload:   mockEvent3.Payload,
			Status:    models.EventProcessing,
			CreatedAt: event3.CreatedAt,
		}

		expected := []models.Event{expectedEvent1, expectedEvent3}

		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		result, err := mockEventRepo.GetEvents(ctx, 10)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCreateEvent(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	mockEvent := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "john"}`),
		Status:  models.EventPending,
	}

	t.Run("Success", func(t *testing.T) {
		err := eventRepo.CreateEvent(ctx, &mockEvent)
		assert.NoError(t, err)

		var result models.Event
		err = db.Where("id = ?", 1).First(&result).Error
		assert.NoError(t, err)

		expected := models.Event{
			ID:        result.ID,
			Name:      mockEvent.Name,
			Payload:   mockEvent.Payload,
			Status:    models.EventPending,
			CreatedAt: result.CreatedAt,
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.CreateEvent(ctx, &models.Event{})

		assert.Error(t, err)
	})
}

func TestMarkProcessing(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	mockEvent := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "john"}`),
		Status:  models.EventPending,
	}
	db.Create(&mockEvent)

	t.Run("Success", func(t *testing.T) {
		err := eventRepo.MarkProcessing(ctx, 1)
		assert.NoError(t, err)

		var result models.Event
		err = db.Where("id = ?", 1).First(&result).Error
		assert.NoError(t, err)

		expected := models.Event{
			ID:        result.ID,
			Name:      mockEvent.Name,
			Payload:   mockEvent.Payload,
			Status:    models.EventProcessing,
			CreatedAt: result.CreatedAt,
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.MarkProcessing(ctx, 1)

		assert.Error(t, err)
	})
}

func TestMarkDone(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	mockEvent := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "john"}`),
		Status:  models.EventPending,
	}
	db.Create(&mockEvent)

	t.Run("Success", func(t *testing.T) {
		err := eventRepo.MarkDone(ctx, 1)
		assert.NoError(t, err)

		var result models.Event
		err = db.Where("id = ?", 1).First(&result).Error
		assert.NoError(t, err)

		expected := models.Event{
			ID:        result.ID,
			Name:      mockEvent.Name,
			Payload:   mockEvent.Payload,
			Status:    models.EventDone,
			CreatedAt: result.CreatedAt,
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.MarkDone(ctx, 1)

		assert.Error(t, err)
	})
}

func TestMarkFailed(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	mockEvent := models.Event{
		Name:    "user.registered",
		Payload: []byte(`{"username": "john"}`),
		Status:  models.EventPending,
	}
	db.Create(&mockEvent)

	t.Run("Success", func(t *testing.T) {
		err := eventRepo.MarkFailed(ctx, 1, "error")
		assert.NoError(t, err)

		var result models.Event
		err = db.Where("id = ?", 1).First(&result).Error
		assert.NoError(t, err)

		expected := models.Event{
			ID:        result.ID,
			Name:      mockEvent.Name,
			Payload:   mockEvent.Payload,
			Status:    models.EventFailed,
			Error:     "error",
			CreatedAt: result.CreatedAt,
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.MarkFailed(ctx, 1, "error")

		assert.Error(t, err)
	})
}
