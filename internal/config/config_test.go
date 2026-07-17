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
			APP_ROOT=/app
			DB_HOST=localhost
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
			REDIS_URL=redis:6379/0
			FASTQC_PATH=/usr/bin/fastqc
			ABRICATE_PATH=/usr/bin/abricate
			MLST_PATH=/usr/bin/mlst
			CHECKM_PATH=/usr/bin/checkm
			KRAKEN2_PATH=/usr/bin/kraken2
			KRAKEN_DB_PATH=/data/kraken_db
			UNICYCLER_PATH=/usr/bin/unicycler
			FASTANI_PATH=/usr/bin/fastani
			SPADES_PATH=/usr/bin/spades
			RESFINDER_DB_PATH=/data/resfinder_db
			POLI_DB_PSEUDO=/dbs/poli/proteins_pseudo_poli.fasta
			POLI_DB_KLEB=/dbs/poli/proteins_kleb_poli.fasta
			POLI_DB_ENTERO=/dbs/poli/proteins_Ecloacae_poli.fasta
			POLI_DB_ACINETO=/dbs/poli/proteins_acineto_poli.fasta
			OTHER_DB_PSEUDO=/dbs/other/proteins_outrasMut_pseudo.fasta
			OTHER_DB_KLEB=/dbs/other/proteins_outrasMut_kleb.fasta
			OTHER_DB_ENTERO=/dbs/other/proteins_outrasMut_Ecloacae.fasta
			OTHER_DB_ACINETO=/dbs/other/proteins_outrasMut_acineto.fasta
			FASTANI_LIST_KLEB=/dbs/fastani/kleb_database/lista-kleb
			FASTANI_LIST_ENTERO=/dbs/fastani/fastANI/list_entero
			FASTANI_LIST_ACINETO=/dbs/fastani/fastANI_acineto/list-acineto
		`
		expectedAppRoot := "/app"
		expectedDbHost := "localhost"
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
		expectedRedisURL := "redis:6379/0"
		expectedFastQCPath := "/usr/bin/fastqc"
		expectedAbricatePath := "/usr/bin/abricate"
		expectedMlstPath := "/usr/bin/mlst"
		expectedCheckMPath := "/usr/bin/checkm"
		expectedKraken2Path := "/usr/bin/kraken2"
		expectedKrakenDBPath := "/data/kraken_db"
		expectedUnicyclerPath := "/usr/bin/unicycler"
		expectedFastaniPath := "/usr/bin/fastani"
		expectedSpadesPath := "/usr/bin/spades"
		expectedResfinderDBPath := "/data/resfinder_db"
		expectedPoliDbPseudo := "/dbs/poli/proteins_pseudo_poli.fasta"
		expectedPoliDbKleb := "/dbs/poli/proteins_kleb_poli.fasta"
		expectedPoliDbEntero := "/dbs/poli/proteins_Ecloacae_poli.fasta"
		expectedPoliDbAcineto := "/dbs/poli/proteins_acineto_poli.fasta"
		expectedOtherDbPseudo := "/dbs/other/proteins_outrasMut_pseudo.fasta"
		expectedOtherDbKleb := "/dbs/other/proteins_outrasMut_kleb.fasta"
		expectedOtherDbEntero := "/dbs/other/proteins_outrasMut_Ecloacae.fasta"
		expectedOtherDbAcineto := "/dbs/other/proteins_outrasMut_acineto.fasta"
		expectedFastaniListKleb := "/dbs/fastani/kleb_database/lista-kleb"
		expectedFastaniListEntero := "/dbs/fastani/fastANI/list_entero"
		expectedFastaniListAcineto := "/dbs/fastani/fastANI_acineto/list-acineto"

		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		testutils.WriteMockEnvFile(t, testEnvFile, envContent)

		err := config.LoadEnvVariables(testEnvFile)
		assert.NoError(t, err)

		assert.Equal(t, expectedAppRoot, os.Getenv("APP_ROOT"), "expected app roots to be equal")
		assert.Equal(t, expectedDbHost, os.Getenv("DB_HOST"), "expected db hosts to be equal")
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
		assert.Equal(t, expectedRedisURL, os.Getenv("REDIS_URL"), "expected redis urls to be equal")
		assert.Equal(t, expectedFastQCPath, os.Getenv("FASTQC_PATH"), "expected fastqc paths to be equal")
		assert.Equal(t, expectedAbricatePath, os.Getenv("ABRICATE_PATH"), "expected abricate paths to be equal")
		assert.Equal(t, expectedMlstPath, os.Getenv("MLST_PATH"), "expected mlst paths to be equal")
		assert.Equal(t, expectedCheckMPath, os.Getenv("CHECKM_PATH"), "expected checkm paths to be equal")
		assert.Equal(t, expectedKraken2Path, os.Getenv("KRAKEN2_PATH"), "expected kraken2 paths to be equal")
		assert.Equal(t, expectedKrakenDBPath, os.Getenv("KRAKEN_DB_PATH"), "expected kraken db paths to be equal")
		assert.Equal(t, expectedUnicyclerPath, os.Getenv("UNICYCLER_PATH"), "expected unicycler paths to be equal")
		assert.Equal(t, expectedFastaniPath, os.Getenv("FASTANI_PATH"), "expected fastani paths to be equal")
		assert.Equal(t, expectedSpadesPath, os.Getenv("SPADES_PATH"), "expected spades paths to be equal")
		assert.Equal(t, expectedResfinderDBPath, os.Getenv("RESFINDER_DB_PATH"), "expected resfinder db paths to be equal")
		assert.Equal(t, expectedPoliDbPseudo, os.Getenv("POLI_DB_PSEUDO"), "expected poli db pseudo to be equal")
		assert.Equal(t, expectedPoliDbKleb, os.Getenv("POLI_DB_KLEB"), "expected poli db kleb to be equal")
		assert.Equal(t, expectedPoliDbEntero, os.Getenv("POLI_DB_ENTERO"), "expected poli db entero to be equal")
		assert.Equal(t, expectedPoliDbAcineto, os.Getenv("POLI_DB_ACINETO"), "expected poli db acineto to be equal")
		assert.Equal(t, expectedOtherDbPseudo, os.Getenv("OTHER_DB_PSEUDO"), "expected other db pseudo to be equal")
		assert.Equal(t, expectedOtherDbKleb, os.Getenv("OTHER_DB_KLEB"), "expected other db kleb to be equal")
		assert.Equal(t, expectedOtherDbEntero, os.Getenv("OTHER_DB_ENTERO"), "expected other db entero to be equal")
		assert.Equal(t, expectedOtherDbAcineto, os.Getenv("OTHER_DB_ACINETO"), "expected other db acineto to be equal")
		assert.Equal(t, expectedFastaniListKleb, os.Getenv("FASTANI_LIST_KLEB"), "expected fastani list kleb to be equal")
	    assert.Equal(t, expectedFastaniListEntero, os.Getenv("FASTANI_LIST_ENTERO"), "expected fastani list entero to be equal")
		assert.Equal(t, expectedFastaniListAcineto, os.Getenv("FASTANI_LIST_ACINETO"), "expected fastani list acineto to be equal")

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

	t.Run("Error - Invalid SMTP port number", func(t *testing.T) {
		os.Unsetenv("PORT")
		os.Unsetenv("SMTP_PORT")
		defer os.Unsetenv("PORT")
		defer os.Unsetenv("SMTP_PORT")

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
			SMTP_PORT=:587
		`
		tempDir := t.TempDir()
		testEnvFile := filepath.Join(tempDir, "test.env")

		testutils.WriteMockEnvFile(t, testEnvFile, envContent)

		err := config.LoadEnvVariables(testEnvFile)
		assert.Error(t, err)
	})
}
