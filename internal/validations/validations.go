package validations

import (
	"errors"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
		models.AdminAnalysisUpdateInput | models.AnalysisTSVDownloadInput
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
