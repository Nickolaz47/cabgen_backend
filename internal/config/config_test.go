package config_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvVariables(t *testing.T) {
	t.Run("Success - Corrected Env file", func(t *testing.T) {
		os.Unsetenv("PORT")
		defer os.Unsetenv("PORT")

		envContent := `
			DB_USER=user
			DB_PASSWORD=password
			DB_NAME=name
			FRONTEND_URL=http://frontend.com
			PORT=8080
			ACCESS_SECRET_KEY=access_secret
			REFRESH_SECRET_KEY=refresh_secret
			ADMIN_PASSWORD=adminpass
			ENVIRONMENT=dev
			API_HOST=localhost
			SENDER_EMAIL=test@gmail.com
			SENDER_PASSWORD=sender_password
			SMTP_HOST=smtp.gmail.com
			SMTP_PORT=587
		`
		expectedUser := "user"
		expectedPassword := "password"
		expectedDbName := "name"
		expectedFrontendUrl := "http://frontend.com"
		expectedPort := 8080
		expectedAccessSecret := "access_secret"
		expectedRefreshSecret := "refresh_secret"
		expectedAdminPassword := "adminpass"
		expectedEnvironment := "dev"
		expectedAPIHost := "localhost"
		expectedSenderEmail := "test@gmail.com"
		expectedSenderPassword := "sender_password"
		expectedSMTPHost := "smtp.gmail.com"
		expectedSMTPPort := 587

		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		testutils.WriteMockEnvFile(t, testEnvFile, envContent)

		err := config.LoadEnvVariables(testEnvFile)
		assert.NoError(t, err)

		assert.Equal(t, expectedUser, os.Getenv("DB_USER"), "expected users to be equal")
		assert.Equal(t, expectedPassword, os.Getenv("DB_PASSWORD"), "expected passwords to be equal")
		assert.Equal(t, expectedDbName, os.Getenv("DB_NAME"), "expected database names to be equal")
		assert.Equal(t, expectedFrontendUrl, os.Getenv("FRONTEND_URL"), "expected URLs to be equal")
		assert.Equal(t, expectedAccessSecret, os.Getenv("ACCESS_SECRET_KEY"), "expected secrets to be equal")
		assert.Equal(t, expectedRefreshSecret, os.Getenv("REFRESH_SECRET_KEY"), "expected secrets to be equal")
		assert.Equal(t, expectedAdminPassword, os.Getenv("ADMIN_PASSWORD"), "expected passwords to be equal")
		assert.Equal(t, expectedEnvironment, os.Getenv("ENVIRONMENT"), "expected environments to be equal")
		assert.Equal(t, expectedAPIHost, os.Getenv("API_HOST"), "expected hosts to be equal")
		assert.Equal(t, expectedSenderEmail, os.Getenv("SENDER_EMAIL"), "expected sender emails to be equal")
		assert.Equal(t, expectedSenderPassword, os.Getenv("SENDER_PASSWORD"), "expected sender passwords to be equal")
		assert.Equal(t, expectedSMTPHost, os.Getenv("SMTP_HOST"), "expected smtp hosts to be equal")

		Port, err := strconv.Atoi(os.Getenv("PORT"))
		assert.NoError(t, err)
		assert.Equal(t, expectedPort, Port, "expected ports to be equal")

		SMTPPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
		assert.NoError(t, err)
		assert.Equal(t, expectedSMTPPort, SMTPPort, "expected ports to be equal")
	})

	t.Run("Error - No default env file", func(t *testing.T) {
		err := config.LoadEnvVariables("")
		assert.Error(t, err)
	})

	t.Run("Error - No customized env file", func(t *testing.T) {
		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		err := config.LoadEnvVariables(testEnvFile)
		assert.Error(t, err)
	})

	t.Run("Error - Invalid port number", func(t *testing.T) {
		os.Unsetenv("PORT")
		defer os.Unsetenv("PORT")

		envContent := `
			DB_USER=user
			DB_PASSWORD=password
			DB_NAME=name
			FRONTEND_URL=http://frontend.com
			PORT=:8080
			ACCESS_SECRET_KEY=access_secret
			REFRESH_SECRET_KEY=refresh_secret
			ADMIN_PASSWORD=adminpass
			ENVIRONMENT=dev
			API_HOST=localhost
			SENDER_EMAIL=test@gmail.com
			SENDER_PASSWORD=sender_password
			SMTP_HOST=smtp.gmail.com
			SMTP_PORT=587
		`
		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		testutils.WriteMockEnvFile(t, testEnvFile, envContent)

		err := config.LoadEnvVariables(testEnvFile)
		assert.Error(t, err)
	})
}
