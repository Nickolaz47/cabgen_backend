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
		result, err := svc.FindAll(context.Background())

		otherOption := models.SelectOption{
			Label: "option.city.other",
			Value: "Other",
		}

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Contains(t, result, otherOption)
	})

	t.Run("Error", func(t *testing.T) {
		origJSON := brazilCitiesJSON
		origCache := brazilCitiesCache

		defer func() {
			brazilCitiesJSON = origJSON
			once = sync.Once{}
			brazilCitiesCache = origCache
		}()

		brazilCitiesJSON = []byte(`{invalid json`)
		once = sync.Once{}
		brazilCitiesCache = nil

		svc := NewCityService()
		result, err := svc.FindAll(context.Background())

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
