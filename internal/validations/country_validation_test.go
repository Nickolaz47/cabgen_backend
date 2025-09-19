package validations_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"

	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
)

func TestValidateCountryCode(t *testing.T) {
	db := testutils.SetupTestRepos()
	mockCountry := testmodels.NewCountry("", "", "", "")
	db.Create(&mockCountry)

	t.Run("Country found", func(t *testing.T) {
		country, ok := validations.ValidateCountryCode("BRA")

		assert.True(t, ok)
		assert.Equal(t, &mockCountry, country)
	})

	t.Run("Country not found", func(t *testing.T) {
		country, ok := validations.ValidateCountryCode("ARG")

		assert.False(t, ok)
		assert.Empty(t, country)
	})
}
