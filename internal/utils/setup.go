package utils

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/security"
	"gorm.io/gorm"
)

func createAdminUser() error {
	var adminUser models.User
	if err := db.DB.Where("username = ?", "admin").First(&adminUser).Error; err == nil {
		return nil
	}

	adminPassword := config.AdminPassword
	if adminPassword == "" {
		return errors.New("admin password is empty")
	}

	hashedPassword, err := security.Hash(adminPassword)
	if err != nil {
		return err
	}

	adminToCreate := models.User{
		Name:        "Cabgen Admin",
		Username:    "admin",
		Email:       "admin@fiocruz.br",
		Password:    hashedPassword,
		CountryCode: "BRA",
		IsActive:    true,
		UserRole:    models.Admin,
		CreatedBy:   "admin",
	}

	if err := db.DB.Create(&adminToCreate).Error; err != nil {
		return fmt.Errorf("cannot create admin user: %v", err)
	}

	return nil
}

func insertCountries() error {
	var count int64
	if err := db.DB.Model(&models.Country{}).Count(&count).Error; err != nil {
		return fmt.Errorf("cannot access countries table: %v", err)
	}

	if count > 0 {
		return nil
	}

	rootDir, err := GetProjectRoot()
	if err != nil {
		return err
	}

	countriesJSON := filepath.Join(rootDir, "internal/data/countries.json")
	data, err := LoadJSONFile[models.Country](countriesJSON)
	if err != nil {
		return err
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			return fmt.Errorf("failed to insert countries data: %v", err)
		}
		return nil
	})
}

func Setup() error {
	if err := insertCountries(); err != nil {
		return err
	}

	if err := createAdminUser(); err != nil {
		return err
	}

	return nil
}
