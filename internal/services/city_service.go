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
	brazilCitiesCache []models.SelectOption
	once              sync.Once
)

type CityService interface {
	FindAll(ctx context.Context) ([]models.SelectOption, error)
}

type cityService struct{}

func NewCityService() CityService {
	return &cityService{}
}

func (s *cityService) FindAll(ctx context.Context) (
	[]models.SelectOption, error) {
	var err error

	once.Do(func() {
		var cities []string
		if unmarshalErr := json.Unmarshal(brazilCitiesJSON,
			&cities); unmarshalErr != nil {
			err = unmarshalErr
			return
		}

		for _, city := range cities {
			brazilCitiesCache = append(brazilCitiesCache, models.SelectOption{
				Label: city,
				Value: city,
			})
		}

		brazilCitiesCache = append(brazilCitiesCache, models.SelectOption{
			Label: "option.city.other",
			Value: "Other",
		})
	})

	if err != nil {
		return nil, err
	}

	return brazilCitiesCache, nil
}
