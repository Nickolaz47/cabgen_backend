package validations_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestApplyUpdateToUser(t *testing.T) {
	user := models.User{
		Name:        "Nicolas",
		Username:    "nick",
		Institution: nil,
		Interest:    nil,
		Role:        nil,
	}

	name := "Nicolas Silva"
	username := "nikol"
	institution := "Fiocruz"
	interest := "Programming"
	role := "Developer"

	updateInput := models.UserUpdateInput{
		Name:        &name,
		Username:    &username,
		Institution: &institution,
		Interest:    &interest,
		Role:        &role,
	}

	validations.ApplyUpdateToUser(&user, &updateInput)

	assert.Equal(t, name, user.Name)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, &institution, user.Institution)
	assert.Equal(t, &interest, user.Interest)
	assert.Equal(t, &role, user.Role)
}
