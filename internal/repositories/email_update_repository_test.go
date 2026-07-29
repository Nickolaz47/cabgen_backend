package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewEmailUpdateRepo(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewEmailUpdateRepo(db)

	assert.NotEmpty(t, repo)
}

func TestEmailUpdateRepositoryCreateRequest(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	repo := repositories.NewEmailUpdateRepo(db)
	userID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		req := testmodels.NewEmailUpdateRequest(userID,
			"old@example.com", "new@example.com",
			"token-create-success", time.Now().Add(1*time.Hour))

		err := repo.CreateRequest(ctx, &req)
		assert.NoError(t, err)

		var dbReq models.EmailUpdateRequest
		err = db.First(&dbReq, "token = ?", "token-create-success").Error
		assert.NoError(t, err)
		assert.Equal(t, req.OldEmail, dbReq.OldEmail)
		assert.Equal(t, req.NewEmail, dbReq.NewEmail)
		assert.Equal(t, userID, dbReq.UserID)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewEmailUpdateRepo(mockDB)
		req := testmodels.NewEmailUpdateRequest(userID,
			"old@example.com", "new@example.com",
			"token-fail", time.Now().Add(1*time.Hour))

		err = mockRepo.CreateRequest(ctx, &req)
		assert.Error(t, err)
	})
}

func TestEmailUpdateRepositoryGetByToken(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	repo := repositories.NewEmailUpdateRepo(db)
	userID := uuid.New()

	req := testmodels.NewEmailUpdateRequest(userID,
		"old@example.com", "new@example.com",
		"token-get-test", time.Now().Add(1*time.Hour).Truncate(time.Second))
	err := db.Create(&req).Error
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		res, err := repo.GetByToken(ctx, "token-get-test")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, req.OldEmail, res.OldEmail)
		assert.Equal(t, req.NewEmail, res.NewEmail)
		assert.Equal(t, req.Token, res.Token)
		assert.WithinDuration(t, req.ExpiresAt, res.ExpiresAt, time.Second)
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		res, err := repo.GetByToken(ctx, "non-existent-token")
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Equal(t, "", res.OldEmail)
	})
}

func TestEmailUpdateRepositoryDeleteRequestsByUserID(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	repo := repositories.NewEmailUpdateRepo(db)
	userID := uuid.New()
	otherUserID := uuid.New()

	req1 := testmodels.NewEmailUpdateRequest(userID,
		"old1@example.com", "new1@example.com",
		"token-del-1", time.Now().Add(1*time.Hour))
	req2 := testmodels.NewEmailUpdateRequest(userID,
		"old2@example.com", "new2@example.com",
		"token-del-2", time.Now().Add(1*time.Hour))
	req3 := testmodels.NewEmailUpdateRequest(otherUserID,
		"old3@example.com", "new3@example.com",
		"token-keep", time.Now().Add(1*time.Hour))

	assert.NoError(t, db.Create(&req1).Error)
	assert.NoError(t, db.Create(&req2).Error)
	assert.NoError(t, db.Create(&req3).Error)

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteRequestsByUserID(ctx, userID)
		assert.NoError(t, err)

		var count int64
		db.Model(&models.EmailUpdateRequest{}).Where("user_id = ?",
			userID).Count(&count)
		assert.Equal(t, int64(0), count)

		var keepCount int64
		db.Model(&models.EmailUpdateRequest{}).Where("user_id = ?",
			otherUserID).Count(&keepCount)
		assert.Equal(t, int64(1), keepCount)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewEmailUpdateRepo(mockDB)
		err = mockRepo.DeleteRequestsByUserID(ctx, userID)
		assert.Error(t, err)
	})
}
