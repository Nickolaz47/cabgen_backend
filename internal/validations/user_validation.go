package validations

import (
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ValidateUpdateInput(c *gin.Context, localizer *i18n.Localizer, input *models.UpdateUserInput) (string, bool) {
	if err := c.ShouldBindJSON(input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			validationErr := ve[0]
			key := "public.auth.register.validation." + validationErr.Field() + "." + validationErr.Tag()
			data := map[string]any{"Param": validationErr.Param()}
			return responses.GetResponseWithData(localizer, key, data), false
		}
		return responses.GetResponse(localizer, responses.RegisterValidationGeneric), false
	}
	return "", true
}

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
