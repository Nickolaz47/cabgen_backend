package validations_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
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
		Email:    &email,
	}

	validations.ApplyAdminUpdateToUser(&user, &updateInput)
	
	assert.Equal(t, name, user.Name)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, &institution, user.Institution)
	assert.Equal(t, &interest, user.Interest)
	assert.Equal(t, &role, user.Role)
}
