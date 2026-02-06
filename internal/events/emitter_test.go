package events_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/events"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewEventEmitter(t *testing.T) {
	mockDB := testutils.NewMockDB()
	repo := repositories.NewEventRepo(mockDB)

	result := events.NewEventEmitter(repo)
	assert.NotNil(t, result)
}

func TestEmit(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewEventRepo(db)
	emitter := events.NewEventEmitter(repo)
	ctx := context.Background()

	correctPayload := `{"username": "john"}`
	wrongPayload := func() {}

	t.Run("Success", func(t *testing.T) {
		err := emitter.Emit(ctx, events.EventUserRegistered, correctPayload)
		assert.NoError(t, err)

		var result models.Event
		err = db.Where("id = ?", 1).First(&result).Error
		assert.NoError(t, err)

		assert.Equal(t, events.EventUserRegistered, result.Name)
	})

	t.Run("Wrong Payload", func(t *testing.T) {
		err := emitter.Emit(ctx, events.EventUserRegistered, wrongPayload)
		assert.Error(t, err)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewEventRepo(mockDB)
		emitter := events.NewEventEmitter(mockRepo)

		err = emitter.Emit(ctx, events.EventUserRegistered, correctPayload)
		assert.Error(t, err)
	})
}
