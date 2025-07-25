package db

import (
	"log"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dns := config.DatabaseConnectionString

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
}

func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Country{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
