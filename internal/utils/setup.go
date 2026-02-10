package utils

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"gorm.io/gorm"
)

func createAdminUser(ctx context.Context, db *gorm.DB) error {
	hasher := security.NewPasswordHasher()

	var adminUser models.User
	if err := db.WithContext(ctx).Where(
		"username = ?", "admin").First(&adminUser).Error; err == nil {
		return nil
	}

	adminPassword := config.AdminPassword
	if adminPassword == "" {
		return errors.New("admin password is empty")
	}

	hashedPassword, err := hasher.Hash(adminPassword)
	if err != nil {
		return err
	}

	countryRepo := repositories.NewCountryRepo(db)
	country, err := countryRepo.GetCountryByCode(ctx, "BRA")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("BRA country not found, did you run the country seed?")
		}
		return fmt.Errorf("cannot fetch country BRA: %w", err)
	}

	adminToCreate := models.User{
		Name:      "Cabgen Admin",
		Username:  "admin",
		Email:     "admin@mail.com",
		Password:  hashedPassword,
		CountryID: country.ID,
		IsActive:  true,
		UserRole:  models.Admin,
		CreatedBy: "admin",
	}

	if err := db.WithContext(ctx).Create(&adminToCreate).Error; err != nil {
		return fmt.Errorf("cannot create admin user: %v", err)
	}

	return nil
}

func insertCountries(ctx context.Context, db *gorm.DB, file string) error {
	repo := repositories.NewCountrySeedRepository(db)

	count, err := repo.Count(ctx)
	if err != nil {
		return fmt.Errorf("cannot access countries table: %w", err)
	}

	if count > 0 {
		return nil
	}

	countries, err := LoadJSONFile[models.Country](file)
	if err != nil {
		return err
	}

	return repo.BulkInsert(ctx, countries)
}

func Setup(db *gorm.DB) error {
	ctx := context.Background()

	rootDir, err := GetProjectRoot()
	if err != nil {
		return err
	}

	countriesJSON := filepath.Join(rootDir, "internal/data/countries.json")
	if err := insertCountries(ctx, db, countriesJSON); err != nil {
		return err
	}

	if err := createAdminUser(ctx, db); err != nil {
		return err
	}

	return nil
}
