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

func TestNewHealthServiceRepo(t *testing.T) {
	db := testutils.NewMockDB()
	healthServiceRepo := repositories.NewHealthServiceRepo(db)

	assert.NotEmpty(t, healthServiceRepo)
}

func TestGetHealthServices(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ := testmodels.NewHealthService(
		uuid.NewString(),
		"Laboratorio Central do Rio de Janeiro",
		models.Public,
		mockCountry,
		"Rio de Janeiro",
		"",
		"",
		"",
		true,
	)
	db.Create(&hServ)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		servs, err := hServicesRepo.GetHealthServices(ctx)
		expected := []models.HealthService{hServ}

		assert.NoError(t, err)
		assert.Equal(t, expected, servs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		labs, err := mockHServRepo.GetHealthServices(ctx)

		assert.Empty(t, labs)
		assert.Error(t, err)
	})
}

func TestGetActiveHealthServices(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ := testmodels.NewHealthService(uuid.NewString(), "HS 1",
		models.Public, mockCountry, "RJ", "", "", "", true)
	hServ2 := testmodels.NewHealthService(uuid.NewString(), "HS 2",
		models.Private, mockCountry, "SP", "", "", "", false)
	db.Create(&hServ)
	db.Create(&hServ2)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		servs, err := hServicesRepo.GetActiveHealthServices(ctx)

		expected := []models.HealthService{hServ}

		assert.NoError(t, err)
		assert.Equal(t, expected, servs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		servs, err := mockHServRepo.GetActiveHealthServices(ctx)

		assert.Empty(t, servs)
		assert.Error(t, err)
	})
}

func TestGetHealthServiceByID(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ := testmodels.NewHealthService(uuid.NewString(), "HS 1",
		models.Public, mockCountry, "RJ", "", "", "", true)
	db.Create(&hServ)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		result, err := hServicesRepo.GetHealthServiceByID(ctx, hServ.ID)

		assert.NoError(t, err)
		assert.Equal(t, hServ, *result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		result, err := mockHServRepo.GetHealthServiceByID(ctx, uuid.New())

		assert.Empty(t, result)
		assert.Error(t, err)
	})
}

func TestGetHealthServicesByName(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ1 := testmodels.NewHealthService(uuid.NewString(),
		"Hospital de Clinicas", models.Public, mockCountry, "RJ", "", "", "",
		true)
	hServ2 := testmodels.NewHealthService(uuid.NewString(),
		"Clinica Sao Jose", models.Private, mockCountry, "SP", "", "", "",
		true)
	db.Create(&hServ1)
	db.Create(&hServ2)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		servs, err := hServicesRepo.GetHealthServicesByName(ctx, "clinica")

		expected := []models.HealthService{hServ1, hServ2}

		assert.NoError(t, err)
		assert.ElementsMatch(t, expected, servs)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		servs, err := mockHServRepo.GetHealthServicesByName(ctx, "Serv")

		assert.Empty(t, servs)
		assert.Error(t, err)
	})
}

func TestGetHealthServiceDuplicate(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	mockHServ := testmodels.NewHealthService(uuid.NewString(),
		"Hospital Central", models.Public, mockCountry, "RJ", "", "", "", true)
	db.Create(&mockHServ)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success - With ID", func(t *testing.T) {
		serv, err := hServicesRepo.GetHealthServiceDuplicate(
			ctx, mockHServ.Name, uuid.New())

		assert.NoError(t, err)
		assert.Equal(t, &mockHServ, serv)
	})

	t.Run("Success - Without ID", func(t *testing.T) {
		serv, err := hServicesRepo.GetHealthServiceDuplicate(
			ctx, mockHServ.Name, uuid.UUID{})

		assert.NoError(t, err)
		assert.Equal(t, &mockHServ, serv)
	})

	t.Run("Error - Record not found", func(t *testing.T) {
		name := "Hospital Falso"
		serv, err := hServicesRepo.GetHealthServiceDuplicate(
			ctx, name, uuid.UUID{})

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, serv)
	})

	t.Run("DB error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		serv, err := mockHServRepo.GetHealthServiceDuplicate(
			ctx, mockHServ.Name, uuid.New())

		assert.Empty(t, serv)
		assert.Error(t, err)
	})
}

func TestCreateHealthService(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ := testmodels.NewHealthService(uuid.NewString(), "HS 1",
		models.Public, mockCountry, "RJ", "", "", "", true)
	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := hServicesRepo.CreateHealthService(ctx, &hServ)
		assert.NoError(t, err)

		var result models.HealthService
		err = db.Where("id = ?", hServ.ID).First(&result).Error

		assert.NoError(t, err)
		assert.Equal(t, hServ.ID, result.ID)
		assert.Equal(t, hServ.Name, result.Name)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		err = mockHServRepo.CreateHealthService(ctx, &hServ)

		assert.Error(t, err)
	})
}

func TestUpdateHealthService(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ := testmodels.NewHealthService(uuid.NewString(), "HS 1",
		models.Public, mockCountry, "RJ", "", "", "", true)
	db.Create(&hServ)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		hServToUpdate := models.HealthService{
			ID:           hServ.ID,
			Name:         "HS Updated",
			Type:         models.Private,
			CountryID:    mockCountry.ID,
			Country:      mockCountry,
			City:         "Sao Paulo",
			Contactant:   "Jose",
			ContactEmail: "jose@email.com",
			ContactPhone: "11999999999",
			IsActive:     true,
		}

		err := hServicesRepo.UpdateHealthService(ctx, &hServToUpdate)
		assert.NoError(t, err)

		var result models.HealthService
		err = db.Where("id = ?", hServ.ID).
			Preload("Country").First(&result).Error

		expected := models.HealthService{
			ID:           hServ.ID,
			Name:         hServToUpdate.Name,
			Type:         hServToUpdate.Type,
			CountryID:    mockCountry.ID,
			Country:      mockCountry,
			City:         hServToUpdate.City,
			Contactant:   hServToUpdate.Contactant,
			ContactEmail: hServToUpdate.ContactEmail,
			ContactPhone: hServToUpdate.ContactPhone,
			IsActive:     hServ.IsActive,
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		err = mockHServRepo.UpdateHealthService(ctx, &hServ)

		assert.Error(t, err)
	})
}

func TestDeleteHealthService(t *testing.T) {
	ctx := context.Background()
	db := testutils.NewMockDB()

	mockCountry := testmodels.NewCountry("", nil)
	db.Create(&mockCountry)

	hServ := testmodels.NewHealthService(uuid.NewString(), "HS 1",
		models.Public, mockCountry, "RJ", "", "", "", true)
	db.Create(&hServ)

	hServicesRepo := repositories.NewHealthServiceRepo(db)

	t.Run("Success", func(t *testing.T) {
		err := hServicesRepo.DeleteHealthService(ctx, &hServ)

		assert.NoError(t, err)

		var result models.HealthService
		err = db.Where("id = ?", hServ.ID).First(&result).Error

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		assert.Empty(t, result)
	})

	t.Run("Error", func(t *testing.T) {
		mockDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		mockHServRepo := repositories.NewHealthServiceRepo(mockDB)
		err = mockHServRepo.DeleteHealthService(ctx, &hServ)

		assert.Error(t, err)
	})
}
