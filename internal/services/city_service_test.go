package services

import (
	"context"
	"sync"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCityFindAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := NewCityService()
		result, err := svc.FindAll(context.Background(), "pt")

		otherOption := models.SelectOption{
			Label: "Outra",
			Value: "Other",
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Contains(t, result, otherOption)
	})

	t.Run("Error", func(t *testing.T) {
		origJSON := brazilCitiesJSON
		origCache := brazilCities

		defer func() {
			brazilCitiesJSON = origJSON
			once = sync.Once{}
			brazilCities = origCache
		}()

		brazilCitiesJSON = []byte(`{invalid json`)
		once = sync.Once{}
		brazilCities = nil

		svc := NewCityService()
		result, err := svc.FindAll(context.Background(), "pt")

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
