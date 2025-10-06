package models_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	t.Run("Success - Slice of bytes", func(t *testing.T) {
		jsonMap := models.JSONMap{}
		value := []byte(`{"en": "Human", "pt": "Humano", "es":"Humano"}`)

		err := jsonMap.Scan(value)

		assert.NoError(t, err)
	})

	t.Run("Success - Empty value", func(t *testing.T) {
		jsonMap := models.JSONMap{}

		err := jsonMap.Scan(nil)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		jsonMap := models.JSONMap{}
		value := 123

		err := jsonMap.Scan(value)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid type for JSONMap:")
	})
}

func TestValue(t *testing.T) {
	t.Run("Success - Slice of bytes", func(t *testing.T) {
		jsonMap := models.JSONMap{"en": "Human"}

		expected := []byte(`{"en":"Human"}`)
		result, err := jsonMap.Value()

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("Success - Empty", func(t *testing.T) {
		jsonMap := models.JSONMap{}

		expected := []byte("{}")
		result, err := jsonMap.Value()

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
