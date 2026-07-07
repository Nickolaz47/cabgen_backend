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
}
