package models

type Country struct {
	ID    uint    `gorm:"primaryKey"`
	Code  string  `gorm:"size:3;uniqueIndex;not null"`
	Names JSONMap `gorm:"type:jsonb;not null"`
	Users []User  `gorm:"foreignKey:CountryID"`
}

type CountryAdminDetailResponse struct {
	Code  string  `json:"code"`
	Names JSONMap `json:"names"`
}

func (c *Country) ToAdminDetailResponse() CountryAdminDetailResponse {
	return CountryAdminDetailResponse{
		Code:  c.Code,
		Names: c.Names,
	}
}

func (c *Country) ToFormResponse(language string) SelectOption {
	if language == "" {
		language = "en"
	}

	return SelectOption{
		Value: c.Code,
		Label: c.Names[language],
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
