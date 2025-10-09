package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	dns := config.DatabaseConnectionString

	newLogger := logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: newLogger,
	})
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

	err := DB.AutoMigrate(&models.User{}, &models.Country{}, &models.Origin{}, &models.Sequencer{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	return nil
}
