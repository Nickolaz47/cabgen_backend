package responses

import "github.com/nicksnyder/go-i18n/v2/i18n"

const (
	HealthMessage                      = "public.health.message"
	RegisterUsernameAlreadyExistsError = "public.auth.register.usernameAlreadyExists.error"
	RegisterEmailAlreadyExistsError    = "public.auth.register.emailAlreadyExists.error"
	GenericInternalServerError         = "generic.internalServer.error"
	RegisterCreateUserError            = "public.auth.register.createUser.error"
	RegisterMessage                    = "public.auth.register.success.message"
	ValidationGeneric                  = "validation.generic"
	RegisterValidationEmailMismatch    = "public.auth.register.validation.emailMismatch"
	RegisterValidationPasswordMismatch = "public.auth.register.validation.passwordMismatch"
	CountryNotFoundError               = "country.notFound.error"
	LoginInvalidCredentialsError       = "validation.invalidCredentials"
	LoginSuccess                       = "public.auth.login.success"
	LoginInactiveUser                  = "public.auth.login.inactiveUser"
	UnauthorizedError                  = "auth.unauthorized"
	TokenExpiredError                  = "auth.tokenExpired"
	LogoutSuccess                      = "public.auth.logout.success"
	TokenRenewed                       = "public.auth.refresh.success"
	UserNotFoundError                  = "user.notFound.error"
	UpdateUserError                    = "user.update.error"
	InvalidUserRoleError               = "admin.user.register.invalidUserRole"
	AdminRegisterSuccess               = "admin.user.register.success"
	UserDeleted                        = "admin.user.delete.success"
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
