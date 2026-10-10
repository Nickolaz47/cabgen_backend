package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/xuri/excelize/v2"
)

const (
	SheetName   = "Samples"
	hiddenSheet = "AcceptedValues"
	maxRows     = 1000
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
	if err := f.SetSheetName("Sheet1", SheetName); err != nil {
		return nil, err
	}

	border := []excelize.Border{
		{Type: "left", Color: "BFBFBF", Style: 1},
		{Type: "right", Color: "BFBFBF", Style: 1},
		{Type: "top", Color: "BFBFBF", Style: 1},
		{Type: "bottom", Color: "BFBFBF", Style: 1},
	}
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:   &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11},
		Fill:   excelize.Fill{Type: "pattern", Color: []string{"0A6354"}, Pattern: 1},
		Border: border,
		Alignment: &excelize.Alignment{
			Horizontal: "center", Vertical: "center", WrapText: true,
		},
	})
	if err != nil {
		return nil, err
	}

	if err := f.SetRowHeight(SheetName, 1, 30); err != nil {
		return nil, err
	}
	for _, col := range columns {
		cell := col.Letter + "1"
		if err := f.SetCellValue(SheetName, cell, col.Header); err != nil {
			return nil, fmt.Errorf("failed to define header %s: %w", col.Header, err)
		}
		if err := f.SetCellStyle(SheetName, cell, cell, headerStyle); err != nil {
			return nil, err
		}
		width := math.Min(math.Max(float64(len(col.Header))+4, 14), 40)
		if err := f.SetColWidth(SheetName, col.Letter, col.Letter, width); err != nil {
			return nil, err
		}
	}

	if err := f.SetPanes(SheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	}); err != nil {
		return nil, err
	}

	if _, err := f.NewSheet(hiddenSheet); err != nil {
		return nil, err
	}

	for _, col := range columns {
		if len(col.LabelValues) == 0 {
			continue
		}

		values := make([]string, len(col.LabelValues))
		copy(values, col.LabelValues)
		sort.Strings(values)

		for i, val := range values {
			cell := fmt.Sprintf("%s%d", col.Letter, i+1)
			if err := f.SetCellValue(hiddenSheet, cell, val); err != nil {
				return nil, err
			}
		}

		dv := excelize.NewDataValidation(true)
		dv.Sqref = fmt.Sprintf("%s2:%s%d", col.Letter, col.Letter, maxRows)
		dv.SetInput("Select or type", "Pick a value from the list or start typing to filter.")
		dv.SetError(excelize.DataValidationErrorStyleStop,
			"Invalid option", "Please select an option from the list.")

		ref := fmt.Sprintf("'%s'!$%s$1:$%s$%d", hiddenSheet, col.Letter, col.Letter, len(values))
		dv.SetSqrefDropList(ref)

		if err := f.AddDataValidation(SheetName, dv); err != nil {
			return nil, fmt.Errorf("failed to add validation for %s: %w", col.Header, err)
		}
	}

	if err := f.SetSheetVisible(hiddenSheet, false); err != nil {
		return nil, err
	}

	idx, err := f.GetSheetIndex(SheetName)
	if err != nil {
		return nil, err
	}
	f.SetActiveSheet(idx)

	return f, nil
}
