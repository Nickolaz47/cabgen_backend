package responses

import "github.com/nicksnyder/go-i18n/v2/i18n"

const (
	HealthMessage                      = "public.health.message"
	RegisterUsernameAlreadyExistsError = "public.auth.register.usernameAlreadyExists.error"
	RegisterEmailAlreadyExistsError    = "public.auth.register.emailAlreadyExists.error"
	GenericInternalServerError         = "generic.internalServer.error"
	RegisterCreateUserError            = "public.auth.register.createUser.error"
	RegisterMessage                    = "public.auth.register.success.message"
	RegisterValidationGeneric          = "public.auth.register.validation.generic"
	RegisterValidationEmailMismatch    = "public.auth.register.validation.emailMismatch"
	RegisterValidationPasswordMismatch = "public.auth.register.validation.passwordMismatch"
)

type APIResponse struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func GetResponse(localizer *i18n.Localizer, messageID string) string {
	msg := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})

	return msg
}

func GetResponseWithData(localizer *i18n.Localizer, messageID string, data map[string]any) string {
	msg := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	return msg
}
