package db_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/db"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/stretchr/testify/assert"
)

func TestNewGormDatabase(t *testing.T) {
	t.Run("Success with sqlite", func(t *testing.T) {
		database, err := db.NewGormDatabase("sqlite", ":memory:")

		assert.NoError(t, err)
		assert.NotNil(t, database)
	})

	t.Run("Unknown driver", func(t *testing.T) {
		database, err := db.NewGormDatabase("unknown", "dsn")

		assert.Error(t, err)
		assert.Empty(t, database)
		assert.Contains(t, err.Error(), "unknown driver")
	})

	t.Run("Postgres retries on unreachable DSN", func(t *testing.T) {
		originalRetries := db.MaxConnectionRetries
		originalDelay := db.ConnectionRetryDelay
		defer func() {
			db.MaxConnectionRetries = originalRetries
			db.ConnectionRetryDelay = originalDelay
		}()

		db.MaxConnectionRetries = 3
		db.ConnectionRetryDelay = time.Millisecond

		database, err := db.NewGormDatabase("postgres",
			"host=localhost user=x password=x dbname=x port=1")

		assert.Error(t, err)
		assert.Nil(t, database)
	})
}

func TestDB(t *testing.T) {
	database, err := db.NewGormDatabase("sqlite", ":memory:")
	assert.NoError(t, err)

	db := database.DB()
	assert.NotEmpty(t, db)
}

func TestMigrate(t *testing.T) {
	models := []any{testmodels.User{}}
	t.Run("Success", func(t *testing.T) {
		database, err := db.NewGormDatabase("sqlite", ":memory:")
		assert.NoError(t, err)

		err = database.Migrate(models...)
		assert.NoError(t, err)
	})

	t.Run("Error - No models", func(t *testing.T) {
		database, err := db.NewGormDatabase("sqlite", ":memory:")
		assert.NoError(t, err)

		err = database.Migrate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "no models provided for migration")
	})

	t.Run("Error - Internal", func(t *testing.T) {
		database, err := db.NewGormDatabase("sqlite", ":memory:")
		assert.NoError(t, err)

		database.Close()

		err = database.Migrate(models...)
		assert.Error(t, err)
	})
}

func TestClose(t *testing.T) {
	database, err := db.NewGormDatabase("sqlite", ":memory:")
	assert.NoError(t, err)

	err = database.Close()
	assert.NoError(t, err)
}
