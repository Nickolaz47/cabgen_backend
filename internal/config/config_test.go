package config_test

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"

	"github.com/stretchr/testify/assert"
)

func writeMockEnvFile(t *testing.T, envFilePath, envContent string) {
	if err := os.WriteFile(envFilePath, []byte(envContent), 0644); err != nil {
		t.Errorf("failed to write mock env file: %v", err)
	}
}

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

		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		writeMockEnvFile(t, testEnvFile, envContent)

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

		SMTPPort, err := strconv.Atoi(os.Getenv("PORT"))

		assert.NoError(t, err)
		assert.Equal(t, expectedPort, SMTPPort, "expected ports to be equal")
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
		`
		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		writeMockEnvFile(t, testEnvFile, envContent)

		err := config.LoadEnvVariables(testEnvFile)
		assert.Error(t, err)
	})
}
