package repositories_test

import (
	"context"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewPasswordResetRepo(t *testing.T) {
	db := testutils.NewMockDB()
	repo := repositories.NewPasswordResetRepo(db)

	assert.NotEmpty(t, repo)
}

func TestPasswordResetRepositoryCreateToken(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	repo := repositories.NewPasswordResetRepo(db)

	t.Run("Success", func(t *testing.T) {
		reset := testmodels.NewPasswordReset("user@example.com",
			"token-create-success", time.Now().Add(1*time.Hour))

		err := repo.CreateToken(ctx, &reset)
		assert.NoError(t, err)

		var dbReset models.PasswordReset
		err = db.First(&dbReset, "token = ?", "token-create-success").Error
		assert.NoError(t, err)
		assert.Equal(t, reset.Email, dbReset.Email)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewPasswordResetRepo(mockDB)
		reset := testmodels.NewPasswordReset("user@example.com",
			"token-fail", time.Now().Add(1*time.Hour))

		err = mockRepo.CreateToken(ctx, &reset)
		assert.Error(t, err)
	})
}

func TestPasswordResetRepositoryGetByToken(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	repo := repositories.NewPasswordResetRepo(db)

	reset := testmodels.NewPasswordReset("user@example.com", "token-get-test",
		time.Now().Add(1*time.Hour).Truncate(time.Second))
	err := db.Create(&reset).Error
	assert.NoError(t, err)

	t.Run("Success", func(t *testing.T) {
		res, err := repo.GetByToken(ctx, "token-get-test")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, reset.Email, res.Email)
		assert.Equal(t, reset.Token, res.Token)
		assert.WithinDuration(t, reset.ExpiresAt, res.ExpiresAt, time.Second)
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		res, err := repo.GetByToken(ctx, "non-existent-token")
		assert.Error(t, err)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Equal(t, "", res.Email)
	})
}

func TestPasswordResetRepositoryDeleteTokensByEmail(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	repo := repositories.NewPasswordResetRepo(db)

	email := "delete-me@example.com"
	reset1 := testmodels.NewPasswordReset(email, "token-del-1",
		time.Now().Add(1*time.Hour))
	reset2 := testmodels.NewPasswordReset(email, "token-del-2",
		time.Now().Add(1*time.Hour))
	reset3 := testmodels.NewPasswordReset("keep-me@example.com", "token-keep",
		time.Now().Add(1*time.Hour))

	assert.NoError(t, db.Create(&reset1).Error)
	assert.NoError(t, db.Create(&reset2).Error)
	assert.NoError(t, db.Create(&reset3).Error)

	t.Run("Success", func(t *testing.T) {
		err := repo.DeleteTokensByEmail(ctx, email)
		assert.NoError(t, err)

		var count int64
		db.Model(&models.PasswordReset{}).Where("email = ?",
			email).Count(&count)
		assert.Equal(t, int64(0), count)

		var keepCount int64
		db.Model(&models.PasswordReset{}).Where("email = ?",
			"keep-me@example.com").Count(&keepCount)
		assert.Equal(t, int64(1), keepCount)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockRepo := repositories.NewPasswordResetRepo(mockDB)
		err = mockRepo.DeleteTokensByEmail(ctx, email)
		assert.Error(t, err)
	})
}
