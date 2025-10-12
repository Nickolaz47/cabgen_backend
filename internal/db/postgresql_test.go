package db_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestConnect(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		err := db.Connect("sqlite", ":memory:")
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := db.Connect("postgres", "invalid-dsn")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to connect")
	})
}

func TestMigrate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		database := testutils.NewMockDB()

		origDB := db.DB
		db.DB = database
		defer func() { db.DB = origDB }()

		err := db.Migrate()
		assert.NoError(t, err)
	})

	t.Run("Empty db", func(t *testing.T) {
		origDB := db.DB
		db.DB = nil
		defer func() { db.DB = origDB }()

		err := db.Migrate()

		assert.Error(t, err)
		assert.ErrorContains(t, err, "DB was not initialized")
	})

	t.Run("Failed to migrate", func(t *testing.T) {
		database, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		origDB := db.DB
		db.DB = database
		defer func() { db.DB = origDB }()

		err = db.Migrate()

		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to migrate database")
	})
}
