package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repository"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewUserRepo(t *testing.T) {
	db := testutils.NewMockDB()
	userRepo := repository.NewUserRepo(db)

	assert.NotEmpty(t, userRepo)
}

func TestGetAllUsers(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()
	filter := models.AdminUserFilter{}

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	mockUser2 := testmodels.NewAdminLoginUser()
	db.Create(&mockUser2)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		users, err := userRepo.GetUsers(ctx, filter)

		for i := range users {
			users[i].CreatedAt = time.Time{}
			users[i].UpdatedAt = time.Time{}
		}
		mockUser.CreatedAt, mockUser.UpdatedAt = time.Time{}, time.Time{}
		mockUser2.CreatedAt, mockUser2.UpdatedAt = time.Time{}, time.Time{}

		expected := []models.User{mockUser, mockUser2}

		assert.NoError(t, err)
		assert.Equal(t, expected, users)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)
		users, err := mockUserRepo.GetUsers(ctx, filter)

		assert.Empty(t, users)
		assert.Error(t, err)
	})
}

func TestGetUserByID(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByID(ctx, mockUser.ID)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByID(ctx, uuid.UUID{})

		assert.Error(t, err)
		assert.Empty(t, user)
	})
}

func TestGetUserByUsername(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByUsername(ctx, mockUser.Username)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByUsername(ctx, "Gohan")

		assert.Error(t, err)
		assert.Empty(t, user)
	})
}

func TestGetUserByEmail(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail(ctx, mockUser.Email)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail(ctx, "i@mail.com")

		assert.Error(t, err)
		assert.Empty(t, user)
	})
}

func TestExistsByEmail(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.ExistsByEmail(
			ctx, &mockUser.Email, uuid.New(),
		)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error - Empty input", func(t *testing.T) {
		user, err := userRepo.ExistsByEmail(
			ctx, nil, uuid.New(),
		)

		assert.Error(t, err)
		assert.Empty(t, user)
		assert.ErrorContains(t, err, "invalid value")
	})

	t.Run("Error - DB", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)
		users, err := mockUserRepo.ExistsByEmail(
			ctx, &mockUser.Email, uuid.New())

		assert.Empty(t, users)
		assert.Error(t, err)
		assert.NotErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestExistsByUsername(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.ExistsByUsername(
			ctx, &mockUser.Username, uuid.New(),
		)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error - Empty input", func(t *testing.T) {
		user, err := userRepo.ExistsByUsername(
			ctx, nil, uuid.New(),
		)

		assert.Error(t, err)
		assert.Empty(t, user)
		assert.ErrorContains(t, err, "invalid value")
	})

	t.Run("Error - DB", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)
		users, err := mockUserRepo.ExistsByUsername(
			ctx, &mockUser.Username, uuid.New())

		assert.Empty(t, users)
		assert.Error(t, err)
		assert.NotErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestCreateUser(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := userRepo.CreateUser(ctx, &mockUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)

		err = mockUserRepo.CreateUser(ctx, &mockUser)

		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := userRepo.UpdateUser(ctx, &mockUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := userRepo.UpdateUser(ctx, &models.User{})

		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	db := testutils.NewMockDB()
	ctx := context.Background()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := userRepo.DeleteUser(ctx, &mockUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := userRepo.DeleteUser(ctx, &models.User{})

		assert.Error(t, err)
	})
}
