package logging

import (
	"context"

	"go.uber.org/zap"
)

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
	EventEmitterError               = "EVENT_EMITTER_ERROR"
	Unauthorized                    = "UNAUTHORIZED"
	CreateFolderError               = "CREATE_FOLDER_ERROR"
	DeleteFolderError               = "DELETE_FOLDER_ERROR"
	DeleteFileError                 = "DELETE_FILE_ERROR"
	MissingFileError                = "MISSING_FILE_ERROR"
	ExceededDownloadLimitError      = "EXCEEDED_DOWNLOAD_LIMIT"
	AsynqTaskError                  = "ASYNQ_TASK_ERROR"
	RedisDispatchError              = "REDIS_DISPATCH_ERROR"
	TicketStatusError               = "TICKET_STATUS_ERROR"
	DeleteActiveTicketError         = "DELETE_ACTIVE_TICKET_ERROR"
	DeletePasswordResetTokenError   = "DELETE_PASSWORD_RESET_TOKEN_ERROR"
	DeleteEmailUpdateRequestError   = "DELETE_EMAIL_UPDATE_REQUEST_ERROR"
	AnalysisRunError                = "ANALYSIS_RUN_ERROR"
)

const (
	TaskEnqueuedSuccess = "TASK_ENQUEUED_SUCCESS"
	EmailSentSuccess    = "EMAIL_SENT_SUCCESS"
)

type requestIDKey struct{}

type userIDKey struct{}

// WithRequestID stores the correlation ID (request ID or task ID) in the
// context so all log lines from the same transaction carry it.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// WithUserID stores the authenticated user's ID in the context so all log
// lines from the same transaction carry it.
func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

// RequestIDFromContext returns the correlation ID stored by WithRequestID.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// UserIDFromContext returns the user ID stored by WithUserID.
func UserIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(userIDKey{}).(string)
	return id
}

func ctxFields(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	var fields []zap.Field
	if id, ok := ctx.Value(requestIDKey{}).(string); ok && id != "" {
		fields = append(fields, zap.String("request_id", id))
	}
	if id, ok := ctx.Value(userIDKey{}).(string); ok && id != "" {
		fields = append(fields, zap.String("user_id", id))
	}

	return fields
}

func ServiceLogging(ctx context.Context, service, function,
	errorType string, err error, extraFields ...zap.Field) []zap.Field {
	fields := append(ctxFields(ctx),
		zap.String("service", service),
		zap.String("func", function),
		zap.String("error_type", errorType),
	)

	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	if len(extraFields) > 0 {
		fields = append(fields, extraFields...)
	}

	return fields
}

func ServiceInfoLogging(ctx context.Context, service, function,
	eventType string, extraFields ...zap.Field) []zap.Field {
	fields := append(ctxFields(ctx),
		zap.String("service", service),
		zap.String("func", function),
		zap.String("event_type", eventType),
	)

	if len(extraFields) > 0 {
		fields = append(fields, extraFields...)
	}

	return fields
}
