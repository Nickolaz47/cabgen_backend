package services

import (
	"context"
	"encoding/json"
	"sync"

	_ "embed"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
)

//go:embed data/brazil_cities.json
var brazilCitiesJSON []byte

var (
	brazilCities []string
	once         sync.Once
)

type CityService interface {
	FindAll(ctx context.Context, language string) (
		[]models.SelectOption, error)
}

type cityService struct{}

func NewCityService() CityService {
	return &cityService{}
}

func (s *cityService) FindAll(ctx context.Context, language string) (
	[]models.SelectOption, error) {
	var err error

	once.Do(func() {
		if unmarshalErr := json.Unmarshal(brazilCitiesJSON,
			&brazilCities); unmarshalErr != nil {
			err = unmarshalErr
		}
	})

	if err != nil {
		return nil, err
	}

	opts := make([]models.SelectOption, 0, len(brazilCities)+1)
	for _, city := range brazilCities {
		opts = append(opts, models.SelectOption{
			Label: city,
			Value: city,
		})
	}

	other, _ := models.ToTranslatedOptionKey("option.city.other", language)
	return append(opts, models.SelectOption{
		Label: other,
		Value: "Other",
	}), nil
}
