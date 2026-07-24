package utils

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
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

type bulkSeeder[T any] interface {
	Count(ctx context.Context) (int64, error)
	BulkInsert(ctx context.Context, items []T) error
}

func seedFromJSON[T any](ctx context.Context, name string,
	repo bulkSeeder[T], file string) error {
	count, err := repo.Count(ctx)
	if err != nil {
		return fmt.Errorf("cannot access seed table: %w", err)
	}
	if count > 0 {
		log.Printf("%s table already populated, skipping seed", name)
		return nil
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Printf("seed file not found, skipping: %s", file)
		return nil
	}

	items, err := LoadJSONFile[T](file)
	if err != nil {
		return err
	}

	if err := repo.BulkInsert(ctx, items); err != nil {
		return err
	}

	log.Printf("seeded %d %s from %s", len(items), name, file)
	return nil
}

func Setup(db *gorm.DB) error {
	ctx := context.Background()

	rootDir, err := GetProjectRoot()
	if err != nil {
		return err
	}

	countriesJSON := filepath.Join(rootDir, "jsons/countries.json")
	if err := seedFromJSON(ctx, "countries",
		repositories.NewCountrySeedRepository(db), countriesJSON); err != nil {
		return err
	}

	microJSON := filepath.Join(rootDir, "jsons/microorganisms.json")
	if err := seedFromJSON(ctx, "microorganisms",
		repositories.NewMicroorganismSeedRepository(db),
		microJSON); err != nil {
		return err
	}

	originsJSON := filepath.Join(rootDir, "jsons/origins.json")
	if err := seedFromJSON(ctx, "origins",
		repositories.NewOriginSeedRepository(db),
		originsJSON); err != nil {
		return err
	}

	sequencersJSON := filepath.Join(rootDir, "jsons/sequencers.json")
	if err := seedFromJSON(ctx, "sequencers",
		repositories.NewSequencerSeedRepository(db),
		sequencersJSON); err != nil {
		return err
	}

	laboratoriesJSON := filepath.Join(rootDir, "jsons/laboratories.json")
	if err := seedFromJSON(ctx, "laboratories",
		repositories.NewLaboratorySeedRepository(db),
		laboratoriesJSON); err != nil {
		return err
	}

	sampleSourcesJSON := filepath.Join(rootDir, "jsons/sample_sources.json")
	if err := seedFromJSON(ctx, "sample_sources",
		repositories.NewSampleSourceSeedRepository(db),
		sampleSourcesJSON); err != nil {
		return err
	}

	if err := createAdminUser(ctx, db); err != nil {
		return err
	}

	return nil
}
