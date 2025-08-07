package validations

import (
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ValidateAdminRegisterInput(c *gin.Context, localizer *i18n.Localizer, newUser *models.AdminRegisterInput) (string, bool) {
	if err := c.ShouldBindJSON(newUser); err != nil {
		var ve validator.ValidationErrors
		var prefix string
		if errors.As(err, &ve) && len(ve) > 0 {
			validationErr := ve[0]
			if validationErr.Field() == "UserRole" {
				prefix = "admin.user.register.validation."
			} else {
				prefix = "public.auth.register.validation."
			}
			key := prefix + validationErr.Field() + "." + validationErr.Tag()
			data := map[string]any{"Param": validationErr.Param()}
			return responses.GetResponseWithData(localizer, key, data), false
		}
		return responses.GetResponse(localizer, responses.RegisterValidationGeneric), false
	}
	return "", true
}
