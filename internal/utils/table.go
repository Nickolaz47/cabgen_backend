package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

var metricsHeaders = []string{
	"coverage", "completeness", "contamination", "genome_size", "n50",
	"primary_species", "secondary_species", "mlst", "poli_mutations",
	"other_mutations", "gene", "resfinder", "vfdb", "plasmid",
}

func GenerateMetricsTSV(analyses []models.AnalysisResponse) ([]byte, error) {
	if len(analyses) == 0 {
		return []byte{}, nil
	}

	buffer := &bytes.Buffer{}
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
			formatTSVValue(r.Coverage),
			r.CheckMCompleteness,
			r.CheckMContamination,
			r.CheckMGenomeSize,
			r.CheckMN50,
			r.PrimarySpeciesName,
			r.SecondarySpeciesName,
			r.MLST,
			strings.Join(r.PoliMutations, ","),
			strings.Join(r.OtherMutations, ","),
			strings.Join(r.ResfinderGenes, ","),
			strings.Join(r.ResfinderBlast, ","),
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
