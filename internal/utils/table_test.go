package utils_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestGenerateMetricsTSV(t *testing.T) {
	t.Run("Success - Empty slice", func(t *testing.T) {
		result, err := utils.GenerateMetricsTSV([]models.AnalysisResponse{})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Success - Single item with all fields", func(t *testing.T) {
		metrics := datatypes.JSON(`{
			"coverage": 30.5,
			"completeness": "95.89",
			"contamination": "1.23",
			"genome_size": "4500000",
			"n50": "12000",
			"primary_species": "Acinetobacter baumannii",
			"secondary_species": "Klebsiella pneumoniae",
			"mlst": "ST502",
			"poli_mutations": ["blaOXA-23", "blaOXA-51"],
			"other_mutations": ["gyrA_S83L"],
			"gene": ["blaOXA-23", "armA"],
			"resfinder": ["blaOXA-23"],
			"vfdb": ["abaum_A"],
			"plasmid": ["IncHI2"]
		}`)
		analyses := []models.AnalysisResponse{{Metrics: metrics}}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		assert.Contains(t, body, "coverage\tcompleteness\tcontamination\tgenome_size\tn50\tprimary_species\tsecondary_species\tmlst\tpoli_mutations\tother_mutations\tgene\tresfinder\tvfdb\tplasmid")
		assert.Contains(t, body, "30.5\t95.89\t1.23\t4500000\t12000\tAcinetobacter baumannii\tKlebsiella pneumoniae\tST502\tblaOXA-23,blaOXA-51\tgyrA_S83L\tblaOXA-23,armA\tblaOXA-23\tabaum_A\tIncHI2")
	})

	t.Run("Success - Single item with empty metrics", func(t *testing.T) {
		analyses := []models.AnalysisResponse{{Metrics: nil}}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		assert.Contains(t, body, "coverage\tcompleteness")
		// Empty row with 14 tab-separated empty cells
		assert.Contains(t, body, "\t\t\t\t\t\t\t\t\t\t\t\t\t\n")
	})

	t.Run("Success - Multiple items", func(t *testing.T) {
		m1 := datatypes.JSON(`{"primary_species": "Species A", "mlst": "ST1"}`)
		m2 := datatypes.JSON(`{"primary_species": "Species B", "mlst": "ST2"}`)
		analyses := []models.AnalysisResponse{
			{Metrics: m1},
			{Metrics: m2},
		}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		assert.Contains(t, body, "Species A")
		assert.Contains(t, body, "Species B")
		assert.Contains(t, body, "ST1")
		assert.Contains(t, body, "ST2")
	})

	t.Run("Success - Array fields joined with comma", func(t *testing.T) {
		metrics := datatypes.JSON(`{
			"gene": ["blaOXA-23", "armA", "blaNDM-1"],
			"poli_mutations": ["mut1"]
		}`)
		analyses := []models.AnalysisResponse{{Metrics: metrics}}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		assert.Contains(t, body, "blaOXA-23,armA,blaNDM-1")
		assert.Contains(t, body, "mut1")
	})

	t.Run("Success - Coverage zero renders empty", func(t *testing.T) {
		metrics := datatypes.JSON(`{"coverage": 0, "primary_species": "Sp"}`)
		analyses := []models.AnalysisResponse{{Metrics: metrics}}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		// First data cell should be empty (coverage=0)
		lines := splitLines(body)
		assert.Len(t, lines, 2) // header + 1 data row
		cells := splitTabs(lines[1])
		assert.Equal(t, "", cells[0]) // coverage empty
		assert.Equal(t, "Sp", cells[5])
	})
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func splitTabs(s string) []string {
	var cells []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\t' {
			cells = append(cells, s[start:i])
			start = i + 1
		}
	}
	cells = append(cells, s[start:])
	return cells
}
