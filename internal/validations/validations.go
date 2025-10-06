package validations

import (
	"errors"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Model interface {
	models.RegisterInput | models.LoginInput |
		models.UpdateUserInput | models.AdminRegisterInput |
		models.AdminUpdateInput | models.OriginCreateInput |
		models.OriginUpdateInput
}

func Validate[T Model](c *gin.Context, localizer *i18n.Localizer, model *T) (string, bool) {
	if err := c.ShouldBindJSON(model); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			validationErr := ve[0]
			key := "validation." + validationErr.Field() + "." + validationErr.Tag()
			data := map[string]any{"Param": validationErr.Param()}
			return responses.GetResponseWithData(localizer, key, data), false
		}
		return responses.GetResponse(localizer, responses.ValidationGeneric), false
	}
	return "", true
}
