package services

import (
	"context"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
)

type SelectOptionsService interface {
	FindAllEnumSelects(ctx context.Context) (*models.EnumSelectsResponse, error)
	FindAllFormSelects(ctx context.Context, language string) (
		*models.FormSelectsResponse, error)
}

type selectOptionsService struct {
	laboratoryRepo    repositories.LaboratoryRepository
	sequencerRepo     repositories.SequencerRepository
	healthServiceRepo repositories.HealthServiceRepository
	originRepo        repositories.OriginRepository
	microorganismRepo repositories.MicroorganismRepository
	sampleSourceRepo  repositories.SampleSourceRepository
}

func NewSelectOptionsService(
	laboratoryRepo repositories.LaboratoryRepository,
	sequencerRepo repositories.SequencerRepository,
	healthServiceRepo repositories.HealthServiceRepository,
	originRepo repositories.OriginRepository,
	microorganismRepo repositories.MicroorganismRepository,
	sampleSourceRepo repositories.SampleSourceRepository,
) SelectOptionsService {
	return &selectOptionsService{
		laboratoryRepo:    laboratoryRepo,
		sequencerRepo:     sequencerRepo,
		healthServiceRepo: healthServiceRepo,
		originRepo:        originRepo,
		microorganismRepo: microorganismRepo,
		sampleSourceRepo:  sampleSourceRepo,
	}
}

func (s *selectOptionsService) FindAllEnumSelects(ctx context.Context) (
	*models.EnumSelectsResponse, error) {
	resp := &models.EnumSelectsResponse{}

	// User Roles
	for _, role := range models.UserRoles {
		resp.Roles = append(resp.Roles, models.SelectOption{
			Label: "option.role." + strings.ToLower(string(role)),
			Value: string(role),
		})
	}

	// Taxons
	for _, taxon := range models.Taxons {
		resp.Taxons = append(resp.Taxons, models.SelectOption{
			Label: "option.taxon." + strings.ToLower(string(taxon)),
			Value: string(taxon),
		})
	}

	// Genders
	for _, gender := range models.Genders {
		resp.Genders = append(resp.Genders, models.SelectOption{
			Label: "option.gender." + strings.ToLower(string(gender)),
			Value: string(gender),
		})
	}

	// Health Service Types
	for _, hsType := range models.HealthServiceTypes {
		resp.HealthServiceTypes = append(resp.HealthServiceTypes,
			models.SelectOption{
				Label: "option.health_service_type." + strings.ToLower(
					string(hsType)),
				Value: string(hsType),
			})
	}

	// Analisis Types
	for _, aType := range models.AnalysisTypes {
		resp.AnalysisTypes = append(resp.AnalysisTypes, models.SelectOption{
			Label: "option.analysis_type." + strings.ToLower(string(aType)),
			Value: string(aType),
		})
	}

	// Languages
	for _, lang := range translation.Languages {
		resp.Languages = append(resp.Languages, models.SelectOption{
			Label: "option.language." + lang,
			Value: lang,
		})
	}

	return resp, nil
}

func (s *selectOptionsService) FindAllFormSelects(ctx context.Context,
	language string) (*models.FormSelectsResponse, error) {
	language = translation.ParseLanguage(language)
	resp := &models.FormSelectsResponse{}

	labs, err := s.laboratoryRepo.GetActiveLaboratories(ctx)
	if err != nil {
		return nil, err
	}
	resp.Laboratories = make([]models.SelectOption, len(labs))
	for i, lab := range labs {
		resp.Laboratories[i] = models.SelectOption{
			Label: lab.Name,
			Value: lab.ID.String(),
		}
	}

	sequencers, err := s.sequencerRepo.GetActiveSequencers(ctx)
	if err != nil {
		return nil, err
	}
	resp.Sequencers = make([]models.SelectOption, len(sequencers))
	for i, seq := range sequencers {
		resp.Sequencers[i] = models.SelectOption{
			Label: seq.Brand,
			Value: seq.ID.String(),
		}
	}

	healthServices, err := s.healthServiceRepo.GetActiveHealthServices(ctx)
	if err != nil {
		return nil, err
	}
	resp.HealthServices = make([]models.SelectOption, len(healthServices))
	for i, hs := range healthServices {
		resp.HealthServices[i] = models.SelectOption{
			Label: hs.Name,
			Value: hs.ID.String(),
		}
	}

	origins, err := s.originRepo.GetActiveOrigins(ctx)
	if err != nil {
		return nil, err
	}
	resp.Origins = make([]models.SelectOption, len(origins))
	for i, origin := range origins {
		resp.Origins[i] = models.SelectOption{
			Label: origin.Names[language],
			Value: origin.ID.String(),
		}
	}

	micros, err := s.microorganismRepo.GetActiveMicroorganisms(ctx)
	if err != nil {
		return nil, err
	}
	resp.Microorganisms = make([]models.SelectOption, len(micros))
	for i, micro := range micros {
		resp.Microorganisms[i] = models.SelectOption{
			Label: micro.Species + " " + micro.Variety[language],
			Value: micro.ID.String(),
		}
	}

	sources, err := s.sampleSourceRepo.GetActiveSampleSources(ctx)
	if err != nil {
		return nil, err
	}
	resp.SampleSources = make([]models.SelectOption, len(sources))
	for i, source := range sources {
		resp.SampleSources[i] = models.SelectOption{
			Label: source.Names[language],
			Value: source.ID.String(),
		}
	}

	return resp, nil
}
