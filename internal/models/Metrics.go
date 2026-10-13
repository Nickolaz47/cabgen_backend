package models

type PublicMetricsResponse struct {
	TotalSamples    int64 `json:"total_samples"`
	TotalCountries  int64 `json:"total_countries"`
	TotalSpecies    int64 `json:"total_species"`
	TotalResistance int64 `json:"total_resistance_genes"`
}

type AdminMetricsResponse struct {
	PublicMetricsResponse
	TotalUsers         int64                `json:"total_users"`
	TotalAnalyses      int64                `json:"total_analyses"`
	UsersByNetwork     NetworkMetric        `json:"users_by_network"`
	SamplesByNetwork   NetworkMetric        `json:"samples_by_network"`
	AnalysesByStatus   AnalysesByStatus     `json:"analyses_by_status"`
	AnalysesByType     []AnalysisTypeMetric `json:"analyses_by_type"`
	TopResistanceGenes []GeneMetric         `json:"top_resistance_genes"`
	SamplesByMonth     []MonthMetric        `json:"samples_by_month"`
	QueueWait          QueueWaitMetric      `json:"queue_wait"`
	TopCountries       []CountryMetric      `json:"top_countries"`
	SpeciesBreakdown   []SpeciesMetric      `json:"species_breakdown"`
}

type AnalysesByStatus struct {
	Done    int64 `json:"done"`
	Running int64 `json:"running"`
	Pending int64 `json:"pending"`
	Failed  int64 `json:"failed"`
}

type NetworkMetric struct {
	InNetwork  int64 `json:"in_network"`
	OutNetwork int64 `json:"out_network"`
}

type GeneMetric struct {
	Gene  string `json:"gene"`
	Count int64  `json:"count"`
}

type MonthMetric struct {
	Month string `json:"month"`
	Count int64  `json:"count"`
}

type QueueWaitMetric struct {
	AvgSeconds float64 `json:"avg_seconds"`
	P95Seconds float64 `json:"p95_seconds"`
}

type AnalysisTypeMetric struct {
	Type               string  `json:"type"`
	Count              int64   `json:"count"`
	AvgDurationSeconds float64 `json:"avg_duration_seconds"`
	P95DurationSeconds float64 `json:"p95_duration_seconds"`
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
