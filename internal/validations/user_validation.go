package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

func ApplyUpdateToUser(user *models.User, input *models.UpdateUserInput) {
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Interest != nil {
		user.Interest = input.Interest
	}
	if input.Role != nil {
		user.Role = input.Role
	}
	if input.Institution != nil {
		user.Institution = input.Institution
	}
}
