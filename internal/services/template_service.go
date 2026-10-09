package services

import (
	"context"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

var templateFreeColumns = map[string]string{
	"A": "origin_code",
	"B": "collection_date",
	"C": "run_number",
	"D": "run_date",
	"G": "date_of_birth",
}

func toMap(options []models.SelectOption) map[string]string {
	m := make(map[string]string)
	for _, opt := range options {
		m[opt.Label] = opt.Value
	}

	return m
}

type TemplateService interface {
	BuildColumns(ctx context.Context, language string) ([]utils.Column, error)
	CreateTemplateTable(ctx context.Context, language string) (
		*excelize.File, error)
}

type templateService struct {
	CountryService       CountryService
	CityService          CityService
	SelectOptionsService SelectOptionsService
	Logger               *zap.Logger
}

func NewTemplateService(countryService CountryService, cityService CityService,
	selectOptionsService SelectOptionsService,
	logger *zap.Logger) TemplateService {
	return &templateService{
		CountryService:       countryService,
		CityService:          cityService,
		SelectOptionsService: selectOptionsService,
		Logger:               logger,
	}
}

func (s *templateService) BuildColumns(ctx context.Context,
	language string) ([]utils.Column, error) {
	cities, err := s.CityService.FindAll(ctx, language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	countries, err := s.CountryService.FindAll(ctx, language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	selectOptions, err := s.SelectOptionsService.FindAllFormSelects(ctx,
		language)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	columns := make([]utils.Column, 0)
	for letter, header := range templateFreeColumns {
		col := utils.NewColumn(letter, header, nil)
		columns = append(columns, col)
	}

	cityColumn := utils.NewColumn("E", "city", toMap(cities))
	genderColumn := utils.NewColumn("F", "gender", toMap(selectOptions.Genders))
	inNetworkOpts := []models.SelectOption{
		{Label: models.ToTranslatedInNetwork(true, language), Value: "true"},
		{Label: models.ToTranslatedInNetwork(false, language), Value: "false"},
	}
	inNetworkColumn := utils.NewColumn("H", "in_network",
		toMap(inNetworkOpts))
	countryColumn := utils.NewColumn("I", "country_code", toMap(countries))
	originColumn := utils.NewColumn("J", "origin_id",
		toMap(selectOptions.Origins))
	sampleSourceColumn := utils.NewColumn("K", "sample_source_id",
		toMap(selectOptions.SampleSources))
	microorganismColumn := utils.NewColumn("L", "microorganism_id",
		toMap(selectOptions.Microorganisms))
	sequencerColumn := utils.NewColumn("M", "sequencer_id",
		toMap(selectOptions.Sequencers))
	laboratoryColumn := utils.NewColumn("N", "laboratory_id",
		toMap(selectOptions.Laboratories))
	healthServiceColumn := utils.NewColumn("O", "health_service_id",
		toMap(selectOptions.HealthServices))

	columns = append(columns, cityColumn, genderColumn, inNetworkColumn,
		countryColumn, originColumn, sampleSourceColumn, microorganismColumn,
		sequencerColumn, laboratoryColumn, healthServiceColumn)

	return columns, nil
}

func (s *templateService) CreateTemplateTable(ctx context.Context,
	language string) (*excelize.File, error) {
	columns, err := s.BuildColumns(ctx, language)
	if err != nil {
		return nil, err
	}

	f, err := utils.GenerateMetadataTemplate(columns)
	if err != nil {
		s.Logger.Error(
			"Service Error", logging.ServiceLogging(ctx,
				"TemplateService", "CreateTemplateTable",
				logging.ExternalRepositoryError, err,
			)...)
		return nil, ErrInternal
	}

	return f, nil
}
