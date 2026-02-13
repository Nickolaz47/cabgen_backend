package validations

import (
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/gin-gonic/gin"
)

func ApplyAdminUpdateToUser(user *models.User, input *models.AdminUserUpdateInput) {
	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	if input.Password != nil {
		user.Password = *input.Password
	}

	if input.UserRole != nil {
		user.UserRole = *input.UserRole
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
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

func ValidateTranslationMap(c *gin.Context, model string, translations map[string]string) (string, bool) {
	localizer := translation.GetLocalizerFromContext(c)
	defaultLanguages := translation.Languages

	var missingLanguage, missingTranslation string
	switch model {
	case "origin":
		missingLanguage, missingTranslation = responses.OriginValidationMissingLanguage, responses.OriginValidationMissingTranslation
	case "sampleSource":
		missingLanguage, missingTranslation = responses.SampleSourceValidationMissingLanguage, responses.SampleSourceValidationMissingTranslation
	case "country":
		missingLanguage, missingTranslation = responses.CountryValidationMissingLanguage, responses.CountryValidationMissingTranslation
	case "microorganism":
		missingLanguage, missingTranslation = responses.MicroorganismValidationMissingLanguage, responses.MicroorganismValidationMissingTranslation
	default:
		missingLanguage, missingTranslation = "", ""
	}

	for _, l := range defaultLanguages {
		value, ok := translations[l]
		if !ok {
			return responses.GetResponseWithData(
				localizer,
				missingLanguage,
				map[string]any{"Param": l},
			), false
		}

		if value == "" {
			return responses.GetResponseWithData(
				localizer,
				missingTranslation,
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

func ApplySequencerUpdate(sequencer *models.Sequencer, input *models.SequencerUpdateInput) {
	if input.Brand != nil {
		sequencer.Brand = *input.Brand
	}

	if input.Model != nil {
		sequencer.Model = *input.Model
	}

	if input.IsActive != nil {
		sequencer.IsActive = *input.IsActive
	}
}

func ApplySampleSourceUpdate(sampleSource *models.SampleSource, input *models.SampleSourceUpdateInput) {
	if input.Names != nil {
		sampleSource.Names = input.Names
	}

	if input.Groups != nil {
		sampleSource.Groups = input.Groups
	}

	if input.IsActive != nil {
		sampleSource.IsActive = *input.IsActive
	}
}

func ApplyLaboratoryUpdate(laboratory *models.Laboratory, input *models.LaboratoryUpdateInput) {
	if input.Name != nil {
		laboratory.Name = *input.Name
	}

	if input.Abbreviation != nil {
		laboratory.Abbreviation = *input.Abbreviation
	}

	if input.IsActive != nil {
		laboratory.IsActive = *input.IsActive
	}
}

func ApplyCountryUpdate(country *models.Country, input *models.CountryUpdateInput) {
	if input.Code != nil {
		country.Code = *input.Code
	}

	if input.Names != nil {
		country.Names = input.Names
	}
}

func ApplyMicroorganismUpdate(microorganism *models.Microorganism,
	input *models.MicroorganismUpdateInput) {
	if input.Taxon != nil {
		microorganism.Taxon = *input.Taxon
	}

	if input.Species != nil {
		microorganism.Species = *input.Species
	}

	if input.Variety != nil {
		microorganism.Variety = input.Variety
	}

	if input.IsActive != nil {
		microorganism.IsActive = *input.IsActive
	}
}
