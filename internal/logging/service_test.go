package logging_test

import (
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

	result := logging.ServiceLogging(service, function, errorType, err)
	expected := []zap.Field{
		zap.String("service", service),
		zap.String("func", function),
		zap.String("error_type", errorType),
		zap.Error(err),
	}

	assert.Equal(t, expected, result)
}
