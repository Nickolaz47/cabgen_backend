package validations

import (
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ValidateLoginInput(c *gin.Context, localizer *i18n.Localizer, newUser *models.LoginInput) (string, bool) {
	if err := c.ShouldBindJSON(newUser); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			validationErr := ve[0]
			key := "public.auth.login.validation." + validationErr.Field() + "." + validationErr.Tag()
			data := map[string]any{"Param": validationErr.Param()}
			return responses.GetResponseWithData(localizer, key, data), false
		}
		return responses.GetResponse(localizer, responses.RegisterValidationGeneric), false
	}
	return "", true
}

func ValidateRegisterInput(c *gin.Context, localizer *i18n.Localizer, newUser *models.RegisterInput) (string, bool) {
	if err := c.ShouldBindJSON(newUser); err != nil {
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

func GetUserTokenFromContext(c *gin.Context) (*models.UserToken, bool) {
	rawUserToken, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userToken, ok := rawUserToken.(*models.UserToken)
	if !ok {
		return nil, false
	}

	return userToken, true
}
