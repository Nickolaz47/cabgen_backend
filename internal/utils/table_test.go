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
			"contigs": "745",
			"primary_species": "Acinetobacter baumannii",
			"secondary_species": "Klebsiella pneumoniae",
			"mlst": "ST502",
			"poli_mutations": ["blaOXA-23", "blaOXA-51"],
			"other_mutations": ["gyrA_S83L"],
		"acquired_resistance": ["blaOXA-23", "armA"],
		"vfdb": ["abaum_A"],
		"plasmid": ["IncHI2"]
		}`)
		analyses := []models.AnalysisResponse{{Metrics: metrics}}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		assert.Contains(t, body, "origin_code\tcoverage\tcompleteness\tcontamination\tgenome_size\tn50\tcontigs\tprimary_species\tsecondary_species\tmlst\tpoli_mutations\tother_mutations\tacquired_resistance\tvfdb\tplasmid")
		assert.Contains(t, body, "\t30.5\t95.89\t1.23\t4500000\t12000\t745\tAcinetobacter baumannii\tKlebsiella pneumoniae\tST502\tblaOXA-23,blaOXA-51\tgyrA_S83L\tblaOXA-23,armA\tabaum_A\tIncHI2")
	})

	t.Run("Success - Single item with empty metrics", func(t *testing.T) {
		analyses := []models.AnalysisResponse{{Metrics: nil}}

		result, err := utils.GenerateMetricsTSV(analyses)

		assert.NoError(t, err)
		body := string(result)
		assert.Contains(t, body, "coverage\tcompleteness")
		assert.Contains(t, body, "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\n")
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
			"acquired_resistance": ["blaOXA-23", "armA", "blaNDM-1"],
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
		lines := splitLines(body)
		assert.Len(t, lines, 2)
		cells := splitTabs(lines[1])
		assert.Equal(t, "", cells[1])
		assert.Equal(t, "Sp", cells[7])
	})
}

func TestNewColumn(t *testing.T) {
	t.Run("Success - Labels sorted without empty values", func(t *testing.T) {
		col := utils.NewColumn("B", "gender", map[string]string{
			"Female":      "F",
			"Male":        "M",
			"Unspecified": "U",
		})

		assert.Equal(t, "B", col.Letter)
		assert.Equal(t, "gender", col.Header)
		assert.Equal(t,
			[]string{"Female", "Male", "Unspecified"}, col.LabelValues)
	})

	t.Run("Success - Empty values", func(t *testing.T) {
		col := utils.NewColumn("A", "city", nil)

		assert.Equal(t, "A", col.Letter)
		assert.Equal(t, "city", col.Header)
		assert.Empty(t, col.LabelValues)
	})
}

func TestColumnFindRealValue(t *testing.T) {
	col := utils.NewColumn("B", "gender", map[string]string{
		"Female":      "F",
		"Male":        "M",
		"Unspecified": "U",
	})

	t.Run("Success - Exact match", func(t *testing.T) {
		real, ok := col.FindRealValue("Male")

		assert.True(t, ok)
		assert.Equal(t, "M", real)
	})

	t.Run("Success - Normalized match", func(t *testing.T) {
		real, ok := col.FindRealValue("  FeMaLE ")

		assert.True(t, ok)
		assert.Equal(t, "F", real)
	})

	t.Run("Success - Not found", func(t *testing.T) {
		real, ok := col.FindRealValue("Other")

		assert.False(t, ok)
		assert.Empty(t, real)
	})
}

func TestGenerateMetadataTemplate(t *testing.T) {
	buildColumns := func() []utils.Column {
		return []utils.Column{
			utils.NewColumn("A", "origin_code", nil),
			utils.NewColumn("B", "gender", map[string]string{
				"Female":      "F",
				"Male":        "M",
				"Unspecified": "U",
			}),
			utils.NewColumn("C", "city", map[string]string{}),
		}
	}

	t.Run("Success - Headers and sheet name", func(t *testing.T) {
		f, err := utils.GenerateMetadataTemplate(buildColumns())

		assert.NoError(t, err)
		assert.Equal(t, "Samples", f.GetSheetName(0))

		for cell, header := range map[string]string{
			"A1": "origin_code",
			"B1": "gender",
			"C1": "city",
		} {
			value, cellErr := f.GetCellValue("Samples", cell)
			assert.NoError(t, cellErr)
			assert.Equal(t, header, value)
		}
	})

	t.Run("Success - Hidden sheet", func(t *testing.T) {
		f, err := utils.GenerateMetadataTemplate(buildColumns())

		assert.NoError(t, err)

		visible, visErr := f.GetSheetVisible("AcceptedValues")
		assert.NoError(t, visErr)
		assert.False(t, visible)

		for cell, expected := range map[string]string{
			"B1": "Female",
			"B2": "Male",
			"B3": "Unspecified",
		} {
			value, cellErr := f.GetCellValue("AcceptedValues", cell)
			assert.NoError(t, cellErr)
			assert.Equal(t, expected, value)
		}
	})

	t.Run("Success - Data validation", func(t *testing.T) {
		f, err := utils.GenerateMetadataTemplate(buildColumns())

		assert.NoError(t, err)

		validations, dvErr := f.GetDataValidations("Samples")
		assert.NoError(t, dvErr)
		if assert.Len(t, validations, 1) {
			assert.Equal(t, "B2:B100", validations[0].Sqref)
			assert.Contains(t, validations[0].Formula1,
				"'AcceptedValues'!$B$1:$B$3")
			if assert.NotNil(t, validations[0].Error) {
				assert.Equal(t, "Please select an option from the list.",
					*validations[0].Error)
			}
		}
	})

	t.Run("Success - Empty columns", func(t *testing.T) {
		f, err := utils.GenerateMetadataTemplate(nil)

		assert.NoError(t, err)

		validations, dvErr := f.GetDataValidations("Samples")
		assert.NoError(t, dvErr)
		assert.Empty(t, validations)
	})

	t.Run("Error - Invalid letter", func(t *testing.T) {
		f, err := utils.GenerateMetadataTemplate([]utils.Column{
			utils.NewColumn("!!", "bad", nil),
		})

		assert.Error(t, err)
		assert.Nil(t, f)
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
