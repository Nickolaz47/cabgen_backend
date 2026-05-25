package repositories_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
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

	mockEvent1 := testmodels.NewEvent(
		uuid.NewString(), "user.registered", []byte(`{"username": "john"}`),
		models.EventPending, "",
	)
	mockEvent2 := testmodels.NewEvent(
		uuid.NewString(), "user.registered", []byte(`{"username": "maria"}`),
		models.EventDone, "",
	)
	mockEvent3 := testmodels.NewEvent(
		uuid.NewString(), "user.registered", []byte(`{"username": "ana"}`),
		models.EventPending, "",
	)

	db.Create(&mockEvent1)
	db.Create(&mockEvent2)
	db.Create(&mockEvent3)

	t.Run("Success", func(t *testing.T) {
		result, err := eventRepo.GetEvents(ctx, 10)
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		var updated1, updated3 models.Event
		assert.NoError(t, db.First(&updated1, "id = ?", mockEvent1.ID).Error)
		assert.NoError(t, db.First(&updated3, "id = ?", mockEvent3.ID).Error)

		assert.Equal(t, models.EventProcessing, updated1.Status)
		assert.Equal(t, models.EventProcessing, updated3.Status)

		ids := []uuid.UUID{result[0].ID, result[1].ID}
		assert.Contains(t, ids, mockEvent1.ID)
		assert.Contains(t, ids, mockEvent3.ID)

		for _, e := range result {
			assert.NotEqual(t, mockEvent2.ID, e.ID)
		}
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

	t.Run("Success", func(t *testing.T) {
		mockEvent := testmodels.NewEvent(
			uuid.NewString(), "user.registered", []byte(`{"username": "john"}`),
			models.EventPending, "",
		)

		err := eventRepo.CreateEvent(ctx, &mockEvent)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, mockEvent.ID)

		var result models.Event
		err = db.First(&result, "id = ?", mockEvent.ID).Error
		assert.NoError(t, err)

		assert.Equal(t, mockEvent.ID, result.ID)
		assert.Equal(t, mockEvent.Name, result.Name)
		assert.Equal(t, mockEvent.Payload, result.Payload)
		assert.Equal(t, models.EventPending, result.Status)
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

	t.Run("Success", func(t *testing.T) {
		mockEvent := testmodels.NewEvent(
			uuid.NewString(), "user.registered", []byte(`{"username": "john"}`),
			models.EventPending, "",
		)

		db.Create(&mockEvent)
		assert.NotEqual(t, uuid.Nil, mockEvent.ID)

		err := eventRepo.MarkProcessing(ctx, mockEvent.ID)
		assert.NoError(t, err)

		var result models.Event
		err = db.First(&result, "id = ?", mockEvent.ID).Error
		assert.NoError(t, err)

		assert.Equal(t, mockEvent.ID, result.ID)
		assert.Equal(t, mockEvent.Name, result.Name)
		assert.Equal(t, mockEvent.Payload, result.Payload)
		assert.Equal(t, models.EventProcessing, result.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.MarkProcessing(ctx, uuid.New())

		assert.Error(t, err)
	})
}

func TestMarkDone(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockEvent := testmodels.NewEvent(
			uuid.NewString(), "user.registered", []byte(`{"username": "john"}`),
			models.EventPending, "",
		)

		db.Create(&mockEvent)
		assert.NotEqual(t, uuid.Nil, mockEvent.ID)

		err := eventRepo.MarkDone(ctx, mockEvent.ID)
		assert.NoError(t, err)

		var result models.Event
		err = db.First(&result, "id = ?", mockEvent.ID).Error
		assert.NoError(t, err)

		assert.Equal(t, mockEvent.ID, result.ID)
		assert.Equal(t, mockEvent.Name, result.Name)
		assert.Equal(t, mockEvent.Payload, result.Payload)
		assert.Equal(t, models.EventDone, result.Status)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.MarkDone(ctx, uuid.New())

		assert.Error(t, err)
	})
}

func TestMarkFailed(t *testing.T) {
	db := testutils.NewMockDB()
	eventRepo := repositories.NewEventRepo(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockEvent := testmodels.NewEvent(
			uuid.NewString(), "user.registered", []byte(`{"username": "john"}`),
			models.EventPending, "",
		)

		db.Create(&mockEvent)
		assert.NotEqual(t, uuid.Nil, mockEvent.ID)

		err := eventRepo.MarkFailed(ctx, mockEvent.ID, "error")
		assert.NoError(t, err)

		var result models.Event
		err = db.First(&result, "id = ?", mockEvent.ID).Error
		assert.NoError(t, err)

		assert.Equal(t, mockEvent.ID, result.ID)
		assert.Equal(t, mockEvent.Name, result.Name)
		assert.Equal(t, mockEvent.Payload, result.Payload)
		assert.Equal(t, models.EventFailed, result.Status)
		assert.Equal(t, "error", result.Error)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockEventRepo := repositories.NewEventRepo(mockDB)
		err = mockEventRepo.MarkFailed(ctx, uuid.New(), "error")

		assert.Error(t, err)
	})
}
