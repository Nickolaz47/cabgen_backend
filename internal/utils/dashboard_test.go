package utils_test

import (
	"strings"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func dashboardRow(t *testing.T, metrics string, sample models.Sample,
	language string) []string {
	t.Helper()

	analysis := models.Analysis{
		Sample:  sample,
		Metrics: datatypes.JSON(metrics),
	}
	result, err := utils.GenerateDashboardTSV([]models.Analysis{analysis},
		language)
	assert.NoError(t, err)

	lines := strings.Split(strings.TrimSuffix(string(result), "\n"), "\n")
	assert.GreaterOrEqual(t, len(lines), 2)
	return strings.Split(lines[1], "\t")
}

func dashboardSample(city string) models.Sample {
	return models.Sample{
		ID:             uuid.New(),
		City:           city,
		CollectionDate: time.Date(2019, time.June, 15, 0, 0, 0, 0, time.UTC),
		SampleSource: models.SampleSource{
			Names: models.JSONMap{
				"pt": "Sangue", "en": "Blood", "es": "Sangre",
			},
		},
	}
}

func TestGenerateDashboardTSV(t *testing.T) {
	t.Run("Success - Empty slice", func(t *testing.T) {
		result, err := utils.GenerateDashboardTSV(nil, "en")

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Success - Full row", func(t *testing.T) {
		metrics := `{
			"primary_species": "Klebsiella pneumoniae subsp pneumoniae",
			"mlst": "klebsiella (ST15)",
			"acquired_resistance": [
				"blaOXA-64_1 (resistance to carbapenem) (allele confidence 99.88)",
				"blaIMP-56_1 (resistance to carbapenem) (allele confidence 88.00)",
				"blaOXA-23_1 (resistance to carbapenem) (allele confidence 100.00)",
				"OqxA_1 (resistance to phenicol/quinolone) (allele confidence 99.23)",
				"sul2_2 (resistance to sulfonamide) (allele confidence 100.00)",
				"dfrA12_1 (resistance to trimethoprim) (allele confidence 99.87)",
				"ant(2'')-Ia_1 (resistance to gentamicin/kanamycin/tobramycin) (allele confidence 100.00)",
				"aac(6')-Ib-cr_1 (resistance to amikacin/kanamycin/quinolone/tobramycin) (allele confidence 99.90)",
				"blaSHV-106_1 (resistance to cephalosporin) (allele confidence 99.88)",
				"blaCTX-M-15_1 (resistance to extended-spectrum cephalosporin) (allele confidence 98.00)"
			],
			"plasmid": [
				"IncFIB(K)_1_Kpn3 (IDENTITY: 98.75 COVERAGE: 100.00 DATABASE: plasmidfinder)",
				"IncN2_1 (IDENTITY: 99.00 COVERAGE: 100.00 DATABASE: plasmidfinder)"
			]
		}`
		sample := dashboardSample("Vitória - ES")

		row := dashboardRow(t, metrics, sample, "pt")

		assert.Equal(t, sample.ID.String(), row[0])
		assert.Equal(t, "Klebsiella pneumoniae subsp pneumoniae", row[1])
		assert.Equal(t,
			"IncFIB(K)_1_Kpn3 IncN2_1", row[2])
		assert.Equal(t, "15", row[3])
		assert.Equal(t, "IMP-56* OXA-23", row[4])
		assert.Equal(t, "OXA-64*", row[5])
		assert.Equal(t, "OqxA", row[6])
		assert.Equal(t, "sul2 dfrA12", row[7])
		assert.Equal(t,
			"ant(2'')-Ia aac(6')-Ib-cr", row[8])
		assert.Equal(t, "SHV-106 CTX-M-15", row[9])
		assert.Equal(t, "Espírito Santo", row[10])
		assert.Equal(t, "Sangue", row[11])
		assert.Equal(t, "2019", row[12])
		assert.Equal(t, "-20.3222", row[13])
		assert.Equal(t, "-40.3381", row[14])
	})

	t.Run("Success - Header row", func(t *testing.T) {
		result, err := utils.GenerateDashboardTSV([]models.Analysis{{
			Sample: dashboardSample("Rio de Janeiro - RJ"),
		}}, "en")

		assert.NoError(t, err)
		firstLine, _, _ := strings.Cut(string(result), "\n")
		assert.Equal(t,
			"ID\tEspécie\tPlasmídeos\tMLST\tCarbapenemases\tOXA-51-like"+
				"\tFluoroquinolonas\tSulfonamidas\tAminoglicosídeos\tESBL"+
				"\tEstado\tMaterial\tData\tLatitude\tLongitude",
			firstLine)
	})

	t.Run("Success - OXA-100 goes to OXA-51-like without star", func(t *testing.T) {
		metrics := `{"acquired_resistance": [
			"blaOXA-100_1 (resistance to carbapenem) (allele confidence 100.00)"
		]}`
		row := dashboardRow(t, metrics, dashboardSample("Salvador - BA"), "en")

		assert.Equal(t, "", row[4])
		assert.Equal(t, "OXA-100", row[5])
	})

	t.Run("Success - Plain quinolone excluded from Fluoroquinolonas",
		func(t *testing.T) {
			metrics := `{"acquired_resistance": [
				"qnrB1_1 (resistance to quinolone) (allele confidence 99.10)"
			]}`
			row := dashboardRow(t, metrics, dashboardSample("Recife - PE"),
				"en")

			assert.Equal(t, "", row[6])
		})

	t.Run("Success - OXA-125 goes to OXA-51-like", func(t *testing.T) {
		metrics := `{"acquired_resistance": [
			"blaOXA-125_1 (resistance to carbapenem) (allele confidence 99.10)"
		]}`
		row := dashboardRow(t, metrics, dashboardSample("Macapá - AP"), "en")

		assert.Equal(t, "OXA-125*", row[5])
	})

	t.Run("Success - MLST without parentheses stays raw", func(t *testing.T) {
		metrics := `{"mlst": "Not available for this species"}`
		row := dashboardRow(t, metrics, dashboardSample("Fortaleza - CE"),
			"en")

		assert.Equal(t, "Not available for this species", row[3])
	})

	t.Run("Success - MLST New ST", func(t *testing.T) {
		metrics := `{"mlst": "efaecium (New ST)"}`
		row := dashboardRow(t, metrics, dashboardSample("Natal - RN"), "en")

		assert.Equal(t, "New ST", row[3])
	})

	t.Run("Success - City without UF leaves state blank", func(t *testing.T) {
		row := dashboardRow(t, `{"primary_species": "Species X"}`,
			dashboardSample("Recife"), "en")

		assert.Equal(t, "", row[10])
		assert.Equal(t, "", row[13])
		assert.Equal(t, "", row[14])
	})

	t.Run("Success - Language picks material name", func(t *testing.T) {
		metrics := `{"primary_species": "Species X"}`
		row := dashboardRow(t, metrics, dashboardSample("Manaus - AM"), "es")

		assert.Equal(t, "Sangre", row[11])
	})

	t.Run("Success - Empty metrics keeps sample fields", func(t *testing.T) {
		sample := dashboardSample("Brasília - DF")
		row := dashboardRow(t, ``, sample, "en")

		assert.Equal(t, sample.ID.String(), row[0])
		assert.Equal(t, "", row[1])
		assert.Equal(t, "", row[4])
		assert.Equal(t, "Distrito Federal", row[10])
		assert.Equal(t, "2019", row[12])
	})
}
