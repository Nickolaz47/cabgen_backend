package logging

import "go.uber.org/zap"

const (
	DatabaseError                 = "DATABASE_ERROR"
	DatabaseNotFoundError         = "DATABASE_NOT_FOUND"
	DatabaseConflictEmailError    = "CONFLICT_EMAIL"
	DatabaseConflictUsernameError = "CONFLICT_USERNAME"
	HasherError                   = "HASHER_ERROR"
	ExternalDatabaseError         = "EXTERNAL_DATABASE_ERROR"
	ExternalDatabaseNotFoundError = "EXTERNAL_DATABASE_NOT_FOUND"
)

func ServiceLogging(service, function, errorType string, err error) []zap.Field {
	return []zap.Field{
		zap.String("service", service),
		zap.String("func", function),
		zap.String("error_type", errorType),
		zap.Error(err),
	}
}
