package db

import (
	"log"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
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
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
}

func Migrate() {
	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	err := DB.AutoMigrate(&models.User{}, &models.Country{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
