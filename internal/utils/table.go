package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/xuri/excelize/v2"
)

var metricsHeaders = []string{
	"origin_code", "coverage", "completeness", "contamination", "genome_size",
	"n50", "contigs", "primary_species", "secondary_species", "mlst", "poli_mutations",
	"other_mutations", "acquired_resistance", "vfdb", "plasmid",
}

func GenerateMetricsTSV(analyses []models.AnalysisResponse) ([]byte, error) {
	if len(analyses) == 0 {
		return []byte{}, nil
	}

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buffer)
	writer.Comma = '\t'

	if err := writer.Write(metricsHeaders); err != nil {
		return nil, err
	}

	for _, a := range analyses {
		var r models.AnalysisResults
		if len(a.Metrics) > 0 {
			_ = json.Unmarshal(a.Metrics, &r)
		}
		row := []string{
			a.Sample,
			formatTSVValue(r.Coverage),
			r.CheckMCompleteness,
			r.CheckMContamination,
			r.CheckMGenomeSize,
			r.CheckMN50,
			r.CheckMContigs,
			r.PrimarySpeciesName,
			r.SecondarySpeciesName,
			r.MLST,
			strings.Join(r.PoliMutations, ","),
			strings.Join(r.OtherMutations, ","),
			strings.Join(r.AcquiredResistance, ","),
			strings.Join(r.VFDB, ","),
			strings.Join(r.PlasmidFinder, ","),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buffer.Bytes(), writer.Error()
}

func formatTSVValue(v any) string {
	switch val := v.(type) {
	case float64:
		if val == 0 {
			return ""
		}
		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

type Column struct {
	Letter        string
	Header        string
	LabelValues   []string
	realValuesMap map[string]string
}

func NewColumn(letter, header string, values map[string]string) Column {
	labelValues := make([]string, 0, len(values))
	realValuesMap := make(map[string]string, len(values))

	if len(values) != 0 {
		for label, real := range values {
			normalizedLabel := strings.TrimSpace(strings.ToLower(label))
			realValuesMap[normalizedLabel] = real
			labelValues = append(labelValues, label)
		}

		slices.Sort(labelValues)
	}

	return Column{
		Letter:        letter,
		Header:        header,
		LabelValues:   labelValues,
		realValuesMap: realValuesMap,
	}
}

func (c *Column) FindRealValue(labelValue string) (string, bool) {
	normalizedLabel := strings.TrimSpace(strings.ToLower(labelValue))
	realValue, ok := c.realValuesMap[normalizedLabel]
	return realValue, ok
}

func GenerateMetadataTemplate(columns []Column) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Samples"
	f.SetSheetName("Sheet1", sheetName)

	for _, col := range columns {
		cellName := fmt.Sprintf("%s1", col.Letter)
		if err := f.SetCellValue(sheetName, cellName, col.Header); err != nil {
			return nil, fmt.Errorf(
				"failed to define metadata header %s: %v", col.Header, err)
		}
	}

	hiddenSheet := "AcceptedValues"
	f.NewSheet(hiddenSheet)
	f.SetSheetVisible(hiddenSheet, false)

	for _, col := range columns {
		if len(col.LabelValues) == 0 {
			continue
		}

		for i, val := range col.LabelValues {
			cellName := fmt.Sprintf("%s%d", col.Letter, i+1)
			f.SetCellValue(hiddenSheet, cellName, val)
		}

		dvOpts := excelize.NewDataValidation(true)
		dvOpts.Sqref = fmt.Sprintf("%s2:%s100", col.Letter, col.Letter)
		dvOpts.SetError(excelize.DataValidationErrorStyleStop,
			"Invalid Option", "Please select an option from the list.")

		formula := fmt.Sprintf("='%s'!$%s$1:$%s$%d",
			hiddenSheet, col.Letter, col.Letter, len(col.LabelValues))
		dvOpts.SetSqrefDropList(formula)

		if err := f.AddDataValidation(sheetName, dvOpts); err != nil {
			return nil, fmt.Errorf("failed to add column validation: %v", err)
		}
	}

	return f, nil
}
