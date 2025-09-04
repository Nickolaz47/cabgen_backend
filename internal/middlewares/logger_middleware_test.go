package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func setupObservedLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)
	return logger, logs
}

func TestLoggerMiddleware(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		consoleLogger, consoleLogs := setupObservedLogger()
		fileLogger, fileLogs := setupObservedLogger()

		w, r := testutils.SetupMiddlewareContext()

		testutils.AddMiddlewares(r, middlewares.LoggerMiddleware(consoleLogger, fileLogger))
		testutils.AddTestGetRoute(r, http.StatusOK)
		testutils.DoGetRequest(r, w)

		assert.Equal(t, http.StatusOK, w.Code)

		assert.True(t, consoleLogs.Len() > 0)
		assert.Contains(t, consoleLogs.All()[0].Message, "Request processed")

		assert.True(t, fileLogs.Len() > 0)
		assert.Equal(t, zapcore.InfoLevel, fileLogs.All()[0].Level)
	})

	t.Run("Client Error", func(t *testing.T) {
		consoleLogger, consoleLogs := setupObservedLogger()
		fileLogger, fileLogs := setupObservedLogger()

		w, r := testutils.SetupMiddlewareContext()

		testutils.AddMiddlewares(r, middlewares.LoggerMiddleware(consoleLogger, fileLogger))
		testutils.AddTestGetRoute(r, http.StatusNotFound)
		testutils.DoGetRequest(r, w)

		assert.Equal(t, http.StatusNotFound, w.Code)

		assert.True(t, consoleLogs.Len() > 0)
		assert.Contains(t, consoleLogs.All()[0].Message, "Client Error")

		assert.True(t, fileLogs.Len() > 0)
		assert.Equal(t, zapcore.WarnLevel, fileLogs.All()[0].Level)
	})

	t.Run("Server Error", func(t *testing.T) {
		consoleLogger, consoleLogs := setupObservedLogger()
		fileLogger, fileLogs := setupObservedLogger()

		w, r := testutils.SetupMiddlewareContext()

		testutils.AddMiddlewares(r, middlewares.LoggerMiddleware(consoleLogger, fileLogger))
		testutils.AddTestGetRoute(r, http.StatusInternalServerError)
		testutils.DoGetRequest(r, w)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		assert.True(t, consoleLogs.Len() > 0)
		assert.Contains(t, consoleLogs.All()[0].Message, "Server Error")

		assert.True(t, fileLogs.Len() > 0)
		assert.Equal(t, zapcore.ErrorLevel, fileLogs.All()[0].Level)
	})
}
