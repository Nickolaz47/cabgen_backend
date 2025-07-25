package responses

import "github.com/nicksnyder/go-i18n/v2/i18n"

const (
	HealthMessage                      = "public.health.message"
	RegisterUsernameAlreadyExistsError = "public.auth.register.usernameAlreadyExists.error"
	RegisterEmailAlreadyExistsError    = "public.auth.register.emailAlreadyExists.error"
	RegisterHashPasswordError          = "public.auth.register.hashPassword.error"
	RegisterCreateUserError            = "public.auth.register.createUser.error"
	RegisterMessage                    = "public.auth.register.success.message"
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
