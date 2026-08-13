package validations

import (
	"errors"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	sanitize "github.com/mrz1836/go-sanitize"
)

type Model interface {
	models.UserRegisterInput | models.LoginInput |
		models.UserUpdateInput | models.AdminUserCreateInput |
		models.AdminUserUpdateInput | models.OriginCreateInput |
		models.OriginUpdateInput | models.SequencerCreateInput |
		models.SequencerUpdateInput | models.SampleSourceCreateInput |
		models.SampleSourceUpdateInput | models.LaboratoryCreateInput |
		models.LaboratoryUpdateInput | models.CountryCreateInput |
		models.CountryUpdateInput | models.MicroorganismCreateInput |
		models.MicroorganismUpdateInput | models.HealthServiceCreateInput |
		models.HealthServiceUpdateInput | models.AdminSampleCreateInput |
		models.AdminSampleUpdateInput | models.SampleCreateInput |
		models.SampleUpdateInput | models.SampleAttachmentInput |
		models.AnalysisCreateInput | models.AdminAnalysisCreateInput |
		models.AdminAnalysisUpdateInput | models.AnalysisTSVDownloadInput |
		models.CreateTicketInput | models.ForgotPasswordInput |
		models.ResetPasswordInput | models.UpdatePasswordInput |
		models.RequestEmailUpdateInput | models.ConfirmEmailUpdateInput
}

func Validate[T Model](
	c *gin.Context, localizer *i18n.Localizer, model *T) (string, bool) {
	if err := c.ShouldBindJSON(model); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) && len(ve) > 0 {
			validationErr := ve[0]
			tag := validationErr.Tag()
			field := validationErr.Field()
			data := map[string]any{"Param": validationErr.Param()}

			namespace := validationErr.StructNamespace()
			structName := strings.SplitN(namespace, ".", 2)[0]

			specificKey := "validation." + structName + "." + field + "." + tag
			if msg := getResponseOrEmpty(localizer, specificKey,
				data); msg != "" {
				return msg, false
			}

			genericKey := "validation." + field + "." + tag
			return responses.GetResponseWithData(localizer, genericKey, data),
				false
		}
		return responses.GetResponse(localizer, responses.ValidationGeneric),
			false
	}

	SanitizeInput(model)

	return validateRequiredDates(model, localizer)
}

func validateRequiredDates[T Model](
	model *T, localizer *i18n.Localizer) (string, bool) {
	switch m := any(*model).(type) {
	case models.AdminSampleCreateInput:
		if m.CollectionDate.IsZero() {
			return responses.GetResponse(localizer,
				"validation.CollectionDate.required"), false
		}
		if m.RunDate.IsZero() {
			return responses.GetResponse(localizer,
				"validation.RunDate.required"), false
		}
	case models.SampleCreateInput:
		if m.CollectionDate.IsZero() {
			return responses.GetResponse(localizer,
				"validation.CollectionDate.required"), false
		}
		if m.RunDate.IsZero() {
			return responses.GetResponse(localizer,
				"validation.RunDate.required"), false
		}
	}
	return "", true
}

func getResponseOrEmpty(localizer *i18n.Localizer, key string,
	data map[string]any) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil || msg == "" {
		return ""
	}
	return msg
}

func SanitizeInput(model any) {
	switch m := model.(type) {
	case *models.UserRegisterInput:
		m.Name = strings.TrimSpace(sanitize.FormalName(m.Name))
		m.Username = strings.ToLower(strings.TrimSpace(m.Username))
		m.Email = sanitize.Email(m.Email, false)
		m.ConfirmEmail = sanitize.Email(m.ConfirmEmail, false)
		m.CountryCode = strings.TrimSpace(m.CountryCode)
		sanitizePtr(m.Interest)
		sanitizePtr(m.Role)
		sanitizePtr(m.Institution)
	case *models.UserUpdateInput:
		sanitizeFormalNamePtr(m.Name)
		sanitizeLowerPtr(m.Username)
		sanitizePtr(m.CountryCode)
		sanitizePtr(m.Interest)
		sanitizePtr(m.Role)
		sanitizePtr(m.Institution)
	case *models.AdminUserCreateInput:
		m.Name = strings.TrimSpace(sanitize.FormalName(m.Name))
		m.Username = strings.ToLower(strings.TrimSpace(m.Username))
		m.Email = sanitize.Email(m.Email, false)
		m.CountryCode = strings.TrimSpace(m.CountryCode)
		sanitizePtr(m.Interest)
		sanitizePtr(m.Role)
		sanitizePtr(m.Institution)
	case *models.AdminUserUpdateInput:
		sanitizeFormalNamePtr(m.Name)
		sanitizeLowerPtr(m.Username)
		sanitizeEmailPtr(m.Email)
		sanitizePtr(m.CountryCode)
		sanitizePtr(m.Interest)
		sanitizePtr(m.Role)
		sanitizePtr(m.Institution)
	case *models.LoginInput:
		m.Username = strings.ToLower(strings.TrimSpace(m.Username))
	case *models.CreateTicketInput:
		m.Name = strings.TrimSpace(sanitize.FormalName(m.Name))
		m.Email = sanitize.Email(m.Email, false)
		m.Institution = strings.TrimSpace(m.Institution)
		m.Subject = sanitize.HTML(sanitize.Scripts(strings.TrimSpace(m.Subject)))
		m.Message = sanitize.HTML(sanitize.Scripts(strings.TrimSpace(m.Message)))
	case *models.AdminSampleCreateInput:
		m.OriginCode = strings.TrimSpace(m.OriginCode)
		m.RunNumber = strings.TrimSpace(m.RunNumber)
		m.CountryCode = strings.TrimSpace(m.CountryCode)
		sanitizePtr(m.City)
	case *models.SampleCreateInput:
		m.OriginCode = strings.TrimSpace(m.OriginCode)
		m.RunNumber = strings.TrimSpace(m.RunNumber)
		m.CountryCode = strings.TrimSpace(m.CountryCode)
		sanitizePtr(m.City)
	case *models.AdminSampleUpdateInput:
		sanitizePtr(m.OriginCode)
		sanitizePtr(m.RunNumber)
		sanitizePtr(m.CountryCode)
		sanitizePtr(m.City)
	case *models.SampleUpdateInput:
		sanitizePtr(m.OriginCode)
		sanitizePtr(m.RunNumber)
		sanitizePtr(m.CountryCode)
		sanitizePtr(m.City)
	case *models.SampleAttachmentInput:
		sanitizePtr(m.Fastq1)
		sanitizePtr(m.Fastq2)
		sanitizePtr(m.Fasta)
	case *models.SequencerCreateInput:
		m.Model = strings.TrimSpace(m.Model)
		m.Brand = strings.TrimSpace(m.Brand)
	case *models.SequencerUpdateInput:
		sanitizePtr(m.Model)
		sanitizePtr(m.Brand)
	case *models.LaboratoryCreateInput:
		m.Name = strings.TrimSpace(m.Name)
		m.Abbreviation = strings.TrimSpace(m.Abbreviation)
	case *models.LaboratoryUpdateInput:
		sanitizePtr(m.Name)
		sanitizePtr(m.Abbreviation)
	case *models.MicroorganismCreateInput:
		m.Species = strings.TrimSpace(m.Species)
		sanitizeMapValues(m.Variety)
	case *models.MicroorganismUpdateInput:
		sanitizePtr(m.Species)
		sanitizeMapValues(m.Variety)
	case *models.HealthServiceCreateInput:
		m.Name = strings.TrimSpace(m.Name)
		m.Type = models.HealthServiceType(strings.TrimSpace(string(m.Type)))
		m.CountryCode = strings.TrimSpace(m.CountryCode)
		sanitizePtr(m.City)
		sanitizePtr(m.Contactant)
		sanitizeEmailPtr(m.ContactEmail)
		sanitizePhonePtr(m.ContactPhone)
	case *models.HealthServiceUpdateInput:
		sanitizePtr(m.Name)
		sanitizeTrimHealthServiceTypePtr(m.Type)
		sanitizePtr(m.CountryCode)
		sanitizePtr(m.City)
		sanitizePtr(m.Contactant)
		sanitizeEmailPtr(m.ContactEmail)
		sanitizePhonePtr(m.ContactPhone)
	case *models.CountryCreateInput:
		m.Code = strings.TrimSpace(m.Code)
		sanitizeMapValues(m.Names)
	case *models.CountryUpdateInput:
		sanitizePtr(m.Code)
		sanitizeMapValues(m.Names)
	case *models.OriginCreateInput:
		sanitizeMapValues(m.Names)
	case *models.OriginUpdateInput:
		sanitizeMapValues(m.Names)
	case *models.SampleSourceCreateInput:
		sanitizeMapValues(m.Names)
		sanitizeMapValues(m.Groups)
	case *models.SampleSourceUpdateInput:
		sanitizeMapValues(m.Names)
		sanitizeMapValues(m.Groups)
	case *models.ForgotPasswordInput:
		m.Email = sanitize.Email(m.Email, false)
	case *models.ResetPasswordInput:
		m.Token = strings.TrimSpace(m.Token)
	case *models.RequestEmailUpdateInput:
		m.NewEmail = sanitize.Email(m.NewEmail, false)
		m.ConfirmNewEmail = sanitize.Email(m.ConfirmNewEmail, false)
	case *models.ConfirmEmailUpdateInput:
		m.Token = strings.TrimSpace(m.Token)
	case *models.AdminAnalysisUpdateInput:
		sanitizePtr(m.FastQC1)
		sanitizePtr(m.FastQC2)
		sanitizePtr(m.ResultsZipPath)
		sanitizePtr(m.ErrorMessage)
	}
}

func sanitizePtr(s *string) {
	if s != nil {
		trimmed := strings.TrimSpace(*s)
		*s = trimmed
	}
}

func sanitizeLowerPtr(s *string) {
	if s != nil {
		*s = strings.ToLower(strings.TrimSpace(*s))
	}
}

func sanitizeFormalNamePtr(s *string) {
	if s != nil {
		*s = strings.TrimSpace(sanitize.FormalName(*s))
	}
}

func sanitizeEmailPtr(s *string) {
	if s != nil {
		*s = sanitize.Email(*s, false)
	}
}

func sanitizePhonePtr(s *string) {
	if s != nil {
		*s = sanitize.PhoneNumber(*s)
	}
}

func sanitizeTrimHealthServiceTypePtr(s *models.HealthServiceType) {
	if s != nil {
		trimmed := strings.TrimSpace(string(*s))
		*s = models.HealthServiceType(trimmed)
	}
}

func sanitizeMapValues(m map[string]string) {
	for k, v := range m {
		m[k] = sanitize.HTML(sanitize.Scripts(strings.TrimSpace(v)))
	}
}
