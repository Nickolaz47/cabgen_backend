package logging_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestSetupLoggersInitializesGlobals(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.log")
	if err := logging.SetupLoggers(tmpFile); err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, logging.ConsoleLogger)
	assert.NotEmpty(t, logging.FileLogger)
	assert.NotEmpty(t, logging.LogFile)
}

func TestFileLoggerWritesLog(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.log")
	if err := logging.SetupLoggers(tmpFile); err != nil {
		t.Fatal(err)
	}

	logging.FileLogger.Info("hello file")
	_ = logging.FileLogger.Sync()

	data, err := os.ReadFile(tmpFile)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "hello file")
}

func TestConsoleLoggerWritesLog(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	tmpFile := filepath.Join(t.TempDir(), "test.log")
	if err := logging.SetupLoggers(tmpFile); err != nil {
		t.Fatal(err)
	}

	logging.ConsoleLogger.Info("hello console")
	_ = logging.ConsoleLogger.Sync()

	w.Close()
	os.Stdout = old

	out, _ := io.ReadAll(r)
	assert.Contains(t, string(out), "hello console")
}

func TestSetupLoggersLogLevel(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectError bool
	}{
		{"Empty defaults to info", "", false},
		{"Debug", "debug", false},
		{"Info", "info", false},
		{"Warn", "warn", false},
		{"Error", "error", false},
		{"Error - Invalid value", "verbose", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LOG_LEVEL", tt.envValue)

			err := logging.SetupLoggers(
				filepath.Join(t.TempDir(), "test.log"))

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, logging.FileLogger)
			}
		})
	}
}
