package models

type PublicMetricsResponse struct {
	TotalSamples    int64 `json:"total_samples"`
	TotalCountries  int64 `json:"total_countries"`
	TotalSpecies    int64 `json:"total_species"`
	TotalResistance int64 `json:"total_resistance_genes"`
}

type AdminMetricsResponse struct {
	PublicMetricsResponse
	TotalUsers       int64            `json:"total_users"`
	TotalAnalyses    int64            `json:"total_analyses"`
	AnalysesByStatus AnalysesByStatus `json:"analyses_by_status"`
	TopCountries     []CountryMetric  `json:"top_countries"`
	SpeciesBreakdown []SpeciesMetric  `json:"species_breakdown"`
}

type AnalysesByStatus struct {
	Done    int64 `json:"done"`
	Running int64 `json:"running"`
	Pending int64 `json:"pending"`
	Failed  int64 `json:"failed"`
}

type CountryMetric struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type SpeciesMetric struct {
	Species string `json:"species"`
	Count   int64  `json:"count"`
}

func (r *AdminMetricsResponse) ToPublicResponse() PublicMetricsResponse {
	return PublicMetricsResponse{
		TotalSamples:    r.TotalSamples,
		TotalCountries:  r.TotalCountries,
		TotalSpecies:    r.TotalSpecies,
		TotalResistance: r.TotalResistance,
	}
}
