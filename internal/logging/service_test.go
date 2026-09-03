package logging_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestServiceLogging(t *testing.T) {
	service := "test"
	function := "testTest"
	errorType := "test failed"
	err := errors.New("test error")

	t.Run("Success", func(t *testing.T) {
		result := logging.ServiceLogging(context.Background(),
			service, function, errorType, err)

		expected := []zap.Field{
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
			zap.Error(err),
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Success - With request_id", func(t *testing.T) {
		ctx := logging.WithRequestID(context.Background(), "req-123")

		result := logging.ServiceLogging(ctx,
			service, function, errorType, err)

		expected := []zap.Field{
			zap.String("request_id", "req-123"),
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
			zap.Error(err),
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Success - With request_id and user_id", func(t *testing.T) {
		ctx := logging.WithUserID(
			logging.WithRequestID(context.Background(), "req-123"),
			"user-456")

		result := logging.ServiceLogging(ctx,
			service, function, errorType, err)

		expected := []zap.Field{
			zap.String("request_id", "req-123"),
			zap.String("user_id", "user-456"),
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
			zap.Error(err),
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Success - With extra fields", func(t *testing.T) {
		result := logging.ServiceLogging(context.Background(),
			service, function, errorType, err,
			zap.String("sample_id", "sample-789"))

		expected := []zap.Field{
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
			zap.Error(err),
			zap.String("sample_id", "sample-789"),
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Error - Nil error omitted", func(t *testing.T) {
		result := logging.ServiceLogging(context.Background(),
			service, function, errorType, nil)

		expected := []zap.Field{
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Success - With user_id only", func(t *testing.T) {
		ctx := logging.WithUserID(context.Background(), "user-456")

		result := logging.ServiceLogging(ctx,
			service, function, errorType, err)

		expected := []zap.Field{
			zap.String("user_id", "user-456"),
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
			zap.Error(err),
		}

		assert.Equal(t, expected, result)
	})
}

func TestLoggingContextGetters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := logging.WithUserID(
			logging.WithRequestID(context.Background(), "req-123"),
			"user-456")

		assert.Equal(t, "req-123",
			logging.RequestIDFromContext(ctx))
		assert.Equal(t, "user-456",
			logging.UserIDFromContext(ctx))
	})

	t.Run("Nil context returns empty", func(t *testing.T) {
		assert.Empty(t, logging.RequestIDFromContext(nil))
		assert.Empty(t, logging.UserIDFromContext(nil))
	})

	t.Run("Empty context returns empty", func(t *testing.T) {
		ctx := context.Background()

		assert.Empty(t, logging.RequestIDFromContext(ctx))
		assert.Empty(t, logging.UserIDFromContext(ctx))
	})
}

func TestServiceInfoLogging(t *testing.T) {
	service := "test"
	function := "testTest"
	eventType := "test event"

	t.Run("Success - With extra fields", func(t *testing.T) {
		ctx := logging.WithRequestID(context.Background(), "req-456")

		result := logging.ServiceInfoLogging(ctx,
			service, function, eventType,
			zap.String("task_type", "email"))

		expected := []zap.Field{
			zap.String("request_id", "req-456"),
			zap.String("service", service),
			zap.String("func", function),
			zap.String("event_type", eventType),
			zap.String("task_type", "email"),
		}

		assert.Equal(t, expected, result)
	})

	t.Run("Success - No extra fields", func(t *testing.T) {
		result := logging.ServiceInfoLogging(context.Background(),
			service, function, eventType)

		expected := []zap.Field{
			zap.String("service", service),
			zap.String("func", function),
			zap.String("event_type", eventType),
		}

		assert.Equal(t, expected, result)
	})
}
