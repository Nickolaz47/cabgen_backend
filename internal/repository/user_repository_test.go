package repository_test

import (
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

func TestGetUsers(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	mockUser2 := testmodels.NewAdminLoginUser()
	db.Create(&mockUser2)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		users, err := userRepo.GetUsers()

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}
		mockUser2.CreatedAt = time.Time{}
		mockUser2.UpdatedAt = time.Time{}
		expected := []models.User{mockUser, mockUser2}

		for i := range users {
			users[i].CreatedAt = time.Time{}
			users[i].UpdatedAt = time.Time{}
		}

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, expected, users)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)
		users, err := mockUserRepo.GetUsers()

		assert.Empty(t, users)
		assert.Error(t, err)
	})
}

func TestGetUserByID(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByID(mockUser.ID)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByID(uuid.UUID{})

		assert.Error(t, err)
		assert.Empty(t, user)
		assert.ErrorContains(t, err, "record not found")
	})
}

func TestGetUserByUsername(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByUsername(mockUser.Username)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByUsername("Gohan")

		assert.Error(t, err)
		assert.Empty(t, user)
		assert.ErrorContains(t, err, "record not found")
	})
}

func TestGetUserByEmail(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail(mockUser.Email)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByEmail("i@mail.com")

		assert.Error(t, err)
		assert.Empty(t, user)
		assert.ErrorContains(t, err, "record not found")
	})
}

func TestGetUserByUsernameOrEmail(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		user, err := userRepo.GetUserByUsernameOrEmail(mockUser.Username, mockUser.Email)

		user.CreatedAt = time.Time{}
		user.UpdatedAt = time.Time{}

		mockUser.CreatedAt = time.Time{}
		mockUser.UpdatedAt = time.Time{}

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
	})

	t.Run("Error", func(t *testing.T) {
		user, err := userRepo.GetUserByUsernameOrEmail("gohan", "g@mail.com")

		assert.Error(t, err)
		assert.Empty(t, user)
		assert.ErrorContains(t, err, "record not found")
	})
}

func TestGetAllAdminUsers(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	mockAdminUser := testmodels.NewAdminLoginUser()
	db.Create(&mockAdminUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		users, err := userRepo.GetAllAdminUsers()

		mockAdminUser.CreatedAt = time.Time{}
		mockAdminUser.UpdatedAt = time.Time{}
		expected := []models.User{mockAdminUser}

		for i := range users {
			users[i].CreatedAt = time.Time{}
			users[i].UpdatedAt = time.Time{}
		}

		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, expected, users)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)

		users, err := mockUserRepo.GetAllAdminUsers()

		assert.Error(t, err)
		assert.Empty(t, users)
	})
}

func TestCreateUser(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := userRepo.CreateUser(&mockUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockUserRepo := repository.NewUserRepo(mockDB)

		err = mockUserRepo.CreateUser(&mockUser)

		assert.Error(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		mockUser.Username = "nikol"

		err := userRepo.UpdateUser(&mockUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := userRepo.UpdateUser(&models.User{})

		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockUser := testmodels.NewLoginUser()
	db.Create(&mockUser)

	userRepo := repository.NewUserRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := userRepo.DeleteUser(&mockUser)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := userRepo.DeleteUser(&models.User{})

		assert.Error(t, err)
	})
}
