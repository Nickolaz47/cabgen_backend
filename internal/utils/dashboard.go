package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"slices"
	"strconv"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

var dashboardHeaders = []string{
	"ID", "Espécie", "Plasmídeos", "MLST", "Carbapenemases",
	"OXA-51-like", "Fluoroquinolonas", "Sulfonamidas",
	"Aminoglicosídeos", "ESBL", "Estado", "Material", "Data",
	"Latitude", "Longitude",
}

var dashboardCarbaCategories = []string{"carbapenem", "xeruborbactam"}

var dashboardFluoroCategories = []string{
	"fluoroquinolone", "phenicol/quinolone",
}

var dashboardSulfoCategories = []string{
	"sulfonamide", "trimetoprim", "trimethoprim",
}

var dashboardAminoCategories = []string{
	"aminoglycoside", "amikacin/kanamycin", "streptomycin",
	"gentamicin", "kanamycin",
}

var dashboardESBLCategories = []string{
	"cephalosporin", "beta-lactam", "(ESBL)",
}

var dashboardOXA51Family = []string{
	"OXA-51", "OXA-64", "OXA-65", "OXA-66", "OXA-69", "OXA-82",
	"OXA-90", "OXA-94", "OXA-98", "OXA-100", "OXA-125", "OXA-132",
	"OXA-259", "OXA-343", "OXA-407", "OXA-408", "OXA-430",
}

type dashboardState struct {
	Nome string
	Lat  string
	Long string
}

var dashboardStates = map[string]dashboardState{
	"AC": {"Acre", "-9.974999", "-67.8243"},
	"AL": {"Alagoas", "-9.6498", "-35.7089"},
	"AP": {"Amapá", "0.034934", "-51.0694"},
	"AM": {"Amazonas", "-3.119", "-60.0217"},
	"BA": {"Bahia", "-12.9704", "-38.5124"},
	"CE": {"Ceará", "-3.7319", "-38.5267"},
	"DF": {"Distrito Federal", "-15.7939", "-47.8828"},
	"ES": {"Espírito Santo", "-20.3222", "-40.3381"},
	"GO": {"Goiás", "-16.6869", "-49.2648"},
	"MA": {"Maranhão", "-2.5307", "-44.3068"},
	"MT": {"Mato Grosso", "-15.6014", "-56.0979"},
	"MS": {"Mato Grosso do Sul", "-20.4697", "-54.6201"},
	"MG": {"Minas Gerais", "-19.9167", "-43.9345"},
	"PA": {"Pará", "-1.4558", "-48.4902"},
	"PB": {"Paraíba", "-7.1195", "-34.8450"},
	"PR": {"Paraná", "-25.4284", "-49.2733"},
	"PE": {"Pernambuco", "-8.0476", "-34.8770"},
	"PI": {"Piauí", "-5.0919", "-42.8034"},
	"RJ": {"Rio de Janeiro", "-22.9035", "-43.2096"},
	"RN": {"Rio Grande do Norte", "-5.7945", "-35.2110"},
	"RS": {"Rio Grande do Sul", "-30.0346", "-51.2177"},
	"RO": {"Rondônia", "-8.7612", "-63.8999"},
	"RR": {"Roraima", "2.8238", "-60.6753"},
	"SC": {"Santa Catarina", "-27.5954", "-48.5480"},
	"SP": {"São Paulo", "-23.5505", "-46.6333"},
	"SE": {"Sergipe", "-10.9472", "-37.0731"},
	"TO": {"Tocantins", "-10.1689", "-48.3317"},
}

func GenerateDashboardTSV(analyses []models.Analysis, language string) (
	[]byte, error) {
	if len(analyses) == 0 {
		return []byte{}, nil
	}

	buffer := &bytes.Buffer{}
	buffer.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buffer)
	writer.Comma = '\t'

	if err := writer.Write(dashboardHeaders); err != nil {
		return nil, err
	}

	for _, a := range analyses {
		var r models.AnalysisResults
		if len(a.Metrics) > 0 {
			_ = json.Unmarshal(a.Metrics, &r)
		}

		state, lat, long := dashboardStateFromCity(a.Sample.City)
		carba, oxa51 := dashboardCarbapenemases(r.AcquiredResistance)

		row := []string{
			a.Sample.ID.String(),
			r.PrimarySpeciesName,
			strings.Join(dashboardPlasmidNames(r.PlasmidFinder), " "),
			dashboardParseMLST(r.MLST),
			strings.Join(carba, " "),
			strings.Join(oxa51, " "),
			strings.Join(dashboardResistanceNames(r.AcquiredResistance,
				dashboardFluoroCategories, false), " "),
			strings.Join(dashboardResistanceNames(r.AcquiredResistance,
				dashboardSulfoCategories, false), " "),
			strings.Join(dashboardResistanceNames(r.AcquiredResistance,
				dashboardAminoCategories, false), " "),
			strings.Join(dashboardResistanceNames(r.AcquiredResistance,
				dashboardESBLCategories, true), " "),
			state,
			a.Sample.SampleSource.Names[language],
			strconv.Itoa(a.Sample.CollectionDate.Year()),
			lat,
			long,
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buffer.Bytes(), writer.Error()
}

func dashboardHasCategory(entry string, categories []string) bool {
	for _, category := range categories {
		if strings.Contains(entry, category) {
			return true
		}
	}
	return false
}

func dashboardRawGene(entry string) string {
	gene, _, _ := strings.Cut(entry, " (")
	return gene
}

func dashboardEntryName(entry string) string {
	name := dashboardRawGene(entry)

	if i := strings.LastIndex(name, "_"); i >= 0 && name[i+1:] != "" &&
		isAllDigits(name[i+1:]) {
		name = name[:i]
	}

	return strings.TrimPrefix(name, "bla")
}

func isAllDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func dashboardEntryConfidence(entry string) (float64, bool) {
	const marker = "(allele confidence "
	_, after, ok := strings.Cut(entry, marker)
	if !ok {
		return 0, false
	}

	rest := after
	rest, _, _ = strings.Cut(rest, ")")
	confidence, err := strconv.ParseFloat(rest, 64)
	if err != nil {
		return 0, false
	}
	return confidence, true
}

func dashboardResistanceNames(entries []string, categories []string,
	blaOnly bool) []string {
	var names []string

	for _, entry := range entries {
		if !dashboardHasCategory(entry, categories) {
			continue
		}
		if blaOnly &&
			!strings.HasPrefix(strings.ToLower(dashboardRawGene(entry)), "bla") {
			continue
		}
		names = append(names, dashboardEntryName(entry))
	}

	return names
}

func dashboardCarbapenemases(entries []string) ([]string, []string) {
	var carbapenemases, oxa51 []string

	for _, entry := range entries {
		if !dashboardHasCategory(entry, dashboardCarbaCategories) {
			continue
		}

		name := dashboardEntryName(entry)
		if confidence, ok := dashboardEntryConfidence(entry); ok &&
			confidence != 100 {
			name += "*"
		}

		if dashboardHasOXA51(name) {
			oxa51 = append(oxa51, name)
			continue
		}
		carbapenemases = append(carbapenemases, name)
	}

	slices.Sort(carbapenemases)
	slices.Sort(oxa51)
	return carbapenemases, oxa51
}

func dashboardHasOXA51(name string) bool {
	for _, member := range dashboardOXA51Family {
		if strings.Contains(name, member) {
			return true
		}
	}
	return false
}

func dashboardPlasmidNames(entries []string) []string {
	var names []string

	for _, entry := range entries {
		name, _, found := strings.Cut(entry, " (IDENTITY:")
		if !found {
			name = entry
		}
		names = append(names, name)
	}

	return names
}

func dashboardParseMLST(mlst string) string {
	open := strings.Index(mlst, "(")
	close := strings.LastIndex(mlst, ")")
	if open < 0 || close < open {
		return mlst
	}

	inner := strings.TrimSpace(mlst[open+1 : close])
	return strings.TrimPrefix(inner, "ST")
}

func dashboardStateFromCity(city string) (string, string, string) {
	i := strings.LastIndex(city, " - ")
	if i < 0 {
		return "", "", ""
	}

	uf := strings.TrimSpace(city[i+len(" - "):])
	state, ok := dashboardStates[uf]
	if !ok {
		return "", "", ""
	}

	return state.Nome, state.Lat, state.Long
}
