package models

type Country struct {
	ID    uint    `gorm:"primaryKey" json:"-"`
	Code  string  `gorm:"size:3;uniqueIndex;not null" json:"code"`
	Names JSONMap `gorm:"type:jsonb;not null" json:"names"`
	Users []User  `gorm:"foreignKey:CountryID" json:"-"`
}

type CountryAdminDetailResponse struct {
	Code  string  `json:"code"`
	Names JSONMap `json:"names"`
}

type CountryFormResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func (c *Country) ToAdminDetailResponse() CountryAdminDetailResponse {
	return CountryAdminDetailResponse{
		Code:  c.Code,
		Names: c.Names,
	}
}

func (c *Country) ToFormResponse(language string) CountryFormResponse {
	if language == "" {
		language = "en"
	}

	return CountryFormResponse{
		Code: c.Code,
		Name: c.Names[language],
	}
}

type CountryCreateInput struct {
	Code  string            `json:"code" binding:"required,len=3"`
	Names map[string]string `json:"names" binding:"required,min=3"`
}

type CountryUpdateInput struct {
	Code  *string           `json:"code,omitempty" binding:"omitempty,len=3"`
	Names map[string]string `json:"names,omitempty" binding:"omitempty,min=3"`
}
