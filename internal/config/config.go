package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DatabaseConnectionString = ""
	AccessKey                = []byte{}
	RefreshKey               = []byte{}
	Port                     = 0
	AdminPassword            = ""
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
		if err := godotenv.Load(envFile); err != nil {
			return err
		}
	} else {
		if err := godotenv.Load(); err != nil {
			return err
		}
	}

	Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return err
	}

	DatabaseConnectionString = fmt.Sprintf(
		"host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=America/Sao_Paulo",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	AccessKey = []byte(os.Getenv("SECRET_ACCESS_KEY"))
	RefreshKey = []byte(os.Getenv("SECRET_REFRESH_KEY"))
	AdminPassword = os.Getenv("ADMIN_PASSWORD")

	return nil
}
