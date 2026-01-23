package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

func ApplyUpdateToUser(user *models.User, input *models.UserUpdateInput) {
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

func IsEmailMatch(email, email2 string) bool {
	return email == email2
}

func IsPasswordMatch(password, password2 string) bool {
	return password == password2
}
