package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(driver, dns string) error {
	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var (
		db  *gorm.DB
		err error
	)

	switch driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
			Logger: newLogger,
		})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dns), &gorm.Config{Logger: newLogger})
	default:
		return fmt.Errorf("unknown driver: %s", driver)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	DB = db
	return nil
}

func Migrate() error {
	if DB == nil {
		return errors.New("DB was not initialized")
	}
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	err := DB.AutoMigrate(
		&models.User{},
		&models.Country{},
		&models.Origin{},
		&models.Sequencer{},
		&models.SampleSource{},
		&models.Laboratory{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}
