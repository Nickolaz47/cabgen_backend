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
	RedisURL                 = ""
	FrontendURL              = ""
	FastQCPath               = ""
	AbricatePath             = ""
	MlstPath                 = ""
	CheckMPath               = ""
	Kraken2Path              = ""
	KrakenDBPath             = ""
	UnicyclerPath            = ""
	FastaniPath              = ""
	SpadesPath               = ""
	ResfinderDBPath          = ""
	PoliDbPseudo             = ""
	PoliDbKleb               = ""
	PoliDbEntero             = ""
	PoliDbAcineto            = ""
	OtherDbPseudo            = ""
	OtherDbKleb              = ""
	OtherDbEntero            = ""
	OtherDbAcineto           = ""
	FastaniListKleb          = ""
	FastaniListEntero        = ""
	FastaniListAcineto       = ""
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
	RedisURL = os.Getenv("REDIS_URL")
	FrontendURL = os.Getenv("FRONTEND_URL")
	FastQCPath = os.Getenv("FASTQC_PATH")
	AbricatePath = os.Getenv("ABRICATE_PATH")
	MlstPath = os.Getenv("MLST_PATH")
	CheckMPath = os.Getenv("CHECKM_PATH")
	Kraken2Path = os.Getenv("KRAKEN2_PATH")
	KrakenDBPath = os.Getenv("KRAKEN_DB_PATH")
	UnicyclerPath = os.Getenv("UNICYCLER_PATH")
	FastaniPath = os.Getenv("FASTANI_PATH")
	SpadesPath = os.Getenv("SPADES_PATH")
	ResfinderDBPath = os.Getenv("RESFINDER_DB_PATH")
	PoliDbPseudo = os.Getenv("POLI_DB_PSEUDO")
	PoliDbKleb = os.Getenv("POLI_DB_KLEB")
	PoliDbEntero = os.Getenv("POLI_DB_ENTERO")
	PoliDbAcineto = os.Getenv("POLI_DB_ACINETO")
	OtherDbPseudo = os.Getenv("OTHER_DB_PSEUDO")
	OtherDbKleb = os.Getenv("OTHER_DB_KLEB")
	OtherDbEntero = os.Getenv("OTHER_DB_ENTERO")
	OtherDbAcineto = os.Getenv("OTHER_DB_ACINETO")
	FastaniListKleb = os.Getenv("FASTANI_LIST_KLEB")
	FastaniListEntero = os.Getenv("FASTANI_LIST_ENTERO")
	FastaniListAcineto = os.Getenv("FASTANI_LIST_ACINETO")

	return nil
}
