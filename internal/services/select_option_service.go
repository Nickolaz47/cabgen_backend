package services

import (
	"context"
	"strings"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type SelectOptionsService interface {
	FindAllEnumSelects(ctx context.Context, language string) (
		*models.EnumSelectsResponse, error)
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

func (s *selectOptionsService) FindAllEnumSelects(ctx context.Context,
	language string) (*models.EnumSelectsResponse, error) {
	resp := &models.EnumSelectsResponse{}
	localizer := i18n.NewLocalizer(translation.Bundle, language)

	// User Roles
	for _, role := range models.UserRoles {
		resp.Roles = append(resp.Roles, models.SelectOption{
			Label: role.ToTranslatedString(language),
			Value: string(role),
		})
	}

	// Taxons
	for _, taxon := range models.Taxons {
		resp.Taxons = append(resp.Taxons, models.SelectOption{
			Label: taxon.ToTranslatedString(language),
			Value: string(taxon),
		})
	}

	// Genders
	for _, gender := range models.Genders {
		resp.Genders = append(resp.Genders, models.SelectOption{
			Label: *gender.ToTranslatedString(language),
			Value: string(gender),
		})
	}

	// Health Service Types
	for _, hsType := range models.HealthServiceTypes {
		resp.HealthServiceTypes = append(resp.HealthServiceTypes,
			models.SelectOption{
				Label: hsType.ToTranslatedString(language),
				Value: string(hsType),
			})
	}

	// Analysis Types
	for _, aType := range models.AnalysisTypes {
		resp.AnalysisTypes = append(resp.AnalysisTypes, models.SelectOption{
			Label: aType.ToTranslatedString(language),
			Value: string(aType),
		})
	}

	// Languages
	for _, lang := range translation.Languages {
		messageID := "option.language." + lang
		label, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: messageID,
		})
		if err != nil {
			label = messageID
		}
		resp.Languages = append(resp.Languages, models.SelectOption{
			Label: label,
			Value: lang,
		})
	}

	return resp, nil
}

func (s *selectOptionsService) FindAllFormSelects(ctx context.Context,
	language string) (*models.FormSelectsResponse, error) {
	language = translation.ParseLanguage(language)
	resp := &models.FormSelectsResponse{}

	translateLabel := func(label string) string {
		if !strings.HasPrefix(label, "option.") {
			return label
		}
		if translated, ok := models.ToTranslatedOptionKey(label, language); ok {
			return translated
		}
		return label
	}

	labs, err := s.laboratoryRepo.GetActiveLaboratories(ctx)
	if err != nil {
		return nil, err
	}
	resp.Laboratories = make([]models.SelectOption, len(labs))
	for i, lab := range labs {
		resp.Laboratories[i] = models.SelectOption{
			Label: translateLabel(lab.Name),
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
			Label: translateLabel(seq.Brand),
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
			Label: translateLabel(hs.Name),
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
			Label: translateLabel(origin.Names[language]),
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
			Label: translateLabel(micro.Species + " " +
				micro.Variety[language]),
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
			Label: translateLabel(source.Names[language]),
			Value: source.ID.String(),
		}
	}

	for _, gender := range models.Genders {
		resp.Genders = append(resp.Genders, models.SelectOption{
			Label: *gender.ToTranslatedString(language),
			Value: string(gender),
		})
	}

	return resp, nil
}
