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

	t.Run("With request_id", func(t *testing.T) {
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

	t.Run("Nil error is omitted", func(t *testing.T) {
		result := logging.ServiceLogging(context.Background(),
			service, function, errorType, nil)

		expected := []zap.Field{
			zap.String("service", service),
			zap.String("func", function),
			zap.String("error_type", errorType),
		}

		assert.Equal(t, expected, result)
	})
}

func TestServiceInfoLogging(t *testing.T) {
	service := "test"
	function := "testTest"
	eventType := "test event"

	t.Run("Success with extra fields", func(t *testing.T) {
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

	t.Run("No extra fields", func(t *testing.T) {
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
