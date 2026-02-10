package logging

import "go.uber.org/zap"

const (
	DatabaseError                   = "DATABASE_ERROR"
	DatabaseNotFoundError           = "DATABASE_NOT_FOUND"
	DatabaseConflictEmailError      = "CONFLICT_EMAIL"
	DatabaseConflictUsernameError   = "CONFLICT_USERNAME"
	HasherError                     = "HASHER_ERROR"
	ExternalRepositoryError         = "EXTERNAL_REPOSITORY_ERROR"
	ExternalRepositoryNotFoundError = "EXTERNAL_REPOSITORY_NOT_FOUND"
	EmailMismatchError              = "EMAIL_MISMATCH"
	PasswordMismatchError           = "PASSWORD_MISMATCH"
	UsernameNotFoundError           = "USERNAME_NOT_FOUND"
	WrongPasswordError              = "WRONG_PASSWORD"
	DisabledUserError               = "DISABLED_USER"
	GetSecretKeyError               = "GET_SECRET_KEY_ERROR"
	GenerateTokenError              = "GENERATE_TOKEN_ERROR"
	ValidateTokenError              = "VALIDATE_TOKEN_ERROR"
	DatabaseConflictError           = "CONFLICT_RECORD"
	SendEmailError                  = "SEND_EMAIL_ERROR"
)

func ServiceLogging(service, function, errorType string, err error) []zap.Field {
	return []zap.Field{
		zap.String("service", service),
		zap.String("func", function),
		zap.String("error_type", errorType),
		zap.Error(err),
	}
}
