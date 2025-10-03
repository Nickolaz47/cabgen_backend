package validations_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestApplyAdminUpdateToUser(t *testing.T) {
	user := models.User{}

	name := "Nicolas Silva"
	username := "nikol"
	email := "nicolas@mail.com"
	institution := "Fiocruz"
	interest := "Programming"
	role := "Developer"

	updateInput := models.AdminUpdateInput{
		UpdateUserInput: models.UpdateUserInput{
			Name:        &name,
			Username:    &username,
			Institution: &institution,
			Interest:    &interest,
			Role:        &role,
		},
		Email: &email,
	}

	validations.ApplyAdminUpdateToUser(&user, &updateInput)

	assert.Equal(t, name, user.Name)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, &institution, user.Institution)
	assert.Equal(t, &interest, user.Interest)
	assert.Equal(t, &role, user.Role)
}

func TestValidateOriginNames(t *testing.T) {
	testutils.SetupTestContext()
	c, _ := testutils.SetupGinContext(
		http.MethodGet, "/", "",
		nil, nil,
	)

	t.Run("Success", func(t *testing.T) {
		originInput := models.OriginCreateInput{
			Names: map[string]string{
				"pt": "Humano",
				"en": "Human",
				"es": "Humano",
			},
			IsActive: true,
		}

		errMsg, ok := validations.ValidateOriginNames(c, &originInput)

		assert.Empty(t, errMsg)
		assert.True(t, ok)
	})

	t.Run("Missing language", func(t *testing.T) {
		originInput := models.OriginCreateInput{
			Names: map[string]string{
				"pt": "Humano",
				"en": "Human",
			},
			IsActive: true,
		}

		errMsg, ok := validations.ValidateOriginNames(c, &originInput)

		assert.Equal(t, errMsg, "Missing es translation.")
		assert.False(t, ok)
	})

	t.Run("Empty translation", func(t *testing.T) {
		originInput := models.OriginCreateInput{
			Names: map[string]string{
				"pt": "",
				"en": "Human",
				"es": "Humano",
			},
			IsActive: true,
		}

		errMsg, ok := validations.ValidateOriginNames(c, &originInput)

		assert.Equal(t, errMsg, "Empty pt translation.")
		assert.False(t, ok)
	})
}

func TestApplyOriginUpdate(t *testing.T) {

}
