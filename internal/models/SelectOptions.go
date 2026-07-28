package models

type SelectOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type EnumSelectsResponse struct {
	Roles              []SelectOption `json:"roles"`
	Taxons             []SelectOption `json:"taxons"`
	Genders            []SelectOption `json:"genders"`
	HealthServiceTypes []SelectOption `json:"health_service_types"`
	AnalysisTypes      []SelectOption `json:"analysis_types"`
	Languages          []SelectOption `json:"languages"`
}

type FormSelectsResponse struct {
	Laboratories   []SelectOption `json:"laboratories"`
	Sequencers     []SelectOption `json:"sequencers"`
	HealthServices []SelectOption `json:"health_services"`
	Origins        []SelectOption `json:"origins"`
	Microorganisms []SelectOption `json:"microorganisms"`
	SampleSources  []SelectOption `json:"sample_sources"`
}
