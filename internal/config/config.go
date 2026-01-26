package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	AppRoot                  = ""
	DatabaseConnectionString = ""
	AccessKey                = []byte{}
	RefreshKey               = []byte{}
	Port                     = 0
	AdminPassword            = ""
	Environment              = ""
	APIHost                  = ""
	SenderEmail              = ""
	SenderPassword           = ""
	SMTPHost                 = ""
	SMTPPort                 = 0
)

/*
LoadEnvVariables loads environment variables from a .env file and assigns them to
global variables.

This function uses the godotenv package to load environment variables from a .env file.
If the .env file cannot be loaded, the function logs a fatal error and terminates the
program.
Otherwise, it assigns the values of specific environment variables
to their corresponding global variables.
*/
func LoadEnvVariables(envFile string) error {
	var err error

	if envFile != "" {
		godotenv.Overload(envFile)
	} else {
		godotenv.Overload()
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		return fmt.Errorf("port variable is missing")
	}

	Port, err = strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	SMTPPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	DatabaseConnectionString = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	AppRoot = os.Getenv("APP_ROOT")
	AccessKey = []byte(os.Getenv("SECRET_ACCESS_KEY"))
	RefreshKey = []byte(os.Getenv("SECRET_REFRESH_KEY"))
	AdminPassword = os.Getenv("ADMIN_PASSWORD")
	Environment = os.Getenv("ENVIRONMENT")
	APIHost = os.Getenv("API_HOST")
	SenderEmail = os.Getenv("SENDER_EMAIL")
	SenderPassword = os.Getenv("SENDER_PASSWORD")
	SMTPHost = os.Getenv("SMTP_HOST")

	return nil
}
