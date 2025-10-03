package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func ApplyAdminUpdateToUser(user *models.User, input *models.AdminUpdateInput) {
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
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

func ValidateOriginNames(c *gin.Context, origin *models.OriginCreateInput) (string, bool) {
	localizer := translation.GetLocalizerFromContext(c)
	defaultLanguages := translation.Languages

	for _, l := range defaultLanguages {
		value, ok := origin.Names[l]
		if !ok {
			return responses.GetResponseWithData(
				localizer,
				responses.OriginValidationMissingLanguage,
				map[string]any{"Param": l},
			), false
		}

		if value == "" {
			return responses.GetResponseWithData(
				localizer,
				responses.OriginValidationMissingTranslation,
				map[string]any{"Param": l},
			), false
		}
	}

	return "", true
}

func ApplyOriginUpdate(origin *models.Origin, input *models.OriginUpdateInput) {
	if input.Names != nil {
		origin.Names = input.Names
	}
	if input.IsActive != nil {
		origin.IsActive = *input.IsActive
	}
}
