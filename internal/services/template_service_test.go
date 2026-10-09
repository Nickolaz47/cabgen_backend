package services_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func templateSuccessMocks() (
	*mocks.MockCityService, *mocks.MockCountryService,
	*mocks.MockSelectOptionsService) {
	citySvc := &mocks.MockCityService{
		FindAllFunc: func(ctx context.Context, language string) (
			[]models.SelectOption, error) {
			return []models.SelectOption{
				{Label: "Rio de Janeiro", Value: "Rio de Janeiro"},
				{Label: "Outro", Value: "Other"},
			}, nil
		},
	}
	countrySvc := &mocks.MockCountryService{
		FindAllFunc: func(ctx context.Context, lang string) (
			[]models.SelectOption, error) {
			return []models.SelectOption{
				{Label: "Brasil", Value: "BRA"},
			}, nil
		},
	}
	originID := uuid.New()
	selectOpts := &mocks.MockSelectOptionsService{
		FindAllFormSelectsFunc: func(ctx context.Context,
			language string) (*models.FormSelectsResponse, error) {
			return &models.FormSelectsResponse{
				Laboratories: []models.SelectOption{
					{Label: "LACEN/RJ", Value: uuid.New().String()},
				},
				Sequencers: []models.SelectOption{
					{Label: "Illumina", Value: uuid.New().String()},
				},
				HealthServices: []models.SelectOption{
					{Label: "Hospital Central", Value: uuid.New().String()},
				},
				Origins: []models.SelectOption{
					{Label: "Humano", Value: originID.String()},
				},
				Microorganisms: []models.SelectOption{
					{Label: "Escherichia coli", Value: uuid.New().String()},
				},
				SampleSources: []models.SelectOption{
					{Label: "Aspirado", Value: uuid.New().String()},
				},
				Genders: []models.SelectOption{
					{Label: "Feminino", Value: "Female"},
					{Label: "Masculino", Value: "Male"},
					{Label: "Não especificado", Value: "Unspecified"},
				},
			}, nil
		},
	}
	return citySvc, countrySvc, selectOpts
}

func columnByLetter(t *testing.T, columns []utils.Column,
	letter string) utils.Column {
	t.Helper()
	for _, col := range columns {
		if col.Letter == letter {
			return col
		}
	}
	t.Fatalf("column %s not found", letter)
	return utils.Column{}
}

func TestTemplateServiceBuildColumns(t *testing.T) {
	t.Run("Success - Full columns", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		var gotLanguage string
		baseFindAll := citySvc.FindAllFunc
		citySvc.FindAllFunc = func(ctx context.Context, language string) (
			[]models.SelectOption, error) {
			gotLanguage = language
			return baseFindAll(ctx, language)
		}

		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)

		columns, err := svc.BuildColumns(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Equal(t, "pt", gotLanguage)
		if assert.Len(t, columns, 15) {
			expectedHeaders := map[string]string{
				"A": "origin_code",
				"B": "collection_date",
				"C": "run_number",
				"D": "run_date",
				"E": "city",
				"F": "gender",
				"G": "date_of_birth",
				"H": "in_network",
				"I": "country_code",
				"J": "origin_id",
				"K": "sample_source_id",
				"L": "microorganism_id",
				"M": "sequencer_id",
				"N": "laboratory_id",
				"O": "health_service_id",
			}
			for letter, header := range expectedHeaders {
				col := columnByLetter(t, columns, letter)
				assert.Equal(t, header, col.Header,
					"header for column %s", letter)
			}

			for _, letter := range []string{"A", "B", "C", "D", "G"} {
				col := columnByLetter(t, columns, letter)
				assert.Empty(t, col.LabelValues,
					"free column %s must not have labels", letter)
			}

			cityCol := columnByLetter(t, columns, "E")
			real, ok := cityCol.FindRealValue("Rio de Janeiro")
			assert.True(t, ok)
			assert.Equal(t, "Rio de Janeiro", real)

			countryCol := columnByLetter(t, columns, "I")
			real, ok = countryCol.FindRealValue("Brasil")
			assert.True(t, ok)
			assert.Equal(t, "BRA", real)

			genderCol := columnByLetter(t, columns, "F")
			real, ok = genderCol.FindRealValue("masculino")
			assert.True(t, ok)
			assert.Equal(t, "Male", real)

			inNetworkCol := columnByLetter(t, columns, "H")
			real, ok = inNetworkCol.FindRealValue("Sim")
			assert.True(t, ok)
			assert.Equal(t, "true", real)
			real, ok = inNetworkCol.FindRealValue("Não")
			assert.True(t, ok)
			assert.Equal(t, "false", real)

			originCol := columnByLetter(t, columns, "J")
			real, ok = originCol.FindRealValue("humano")
			assert.True(t, ok)
			assert.NotEmpty(t, real)
		}
	})

	t.Run("Success - Empty lists", func(t *testing.T) {
		citySvc := &mocks.MockCityService{
			FindAllFunc: func(ctx context.Context, language string) (
				[]models.SelectOption, error) {
				return []models.SelectOption{}, nil
			},
		}
		countrySvc := &mocks.MockCountryService{
			FindAllFunc: func(ctx context.Context, lang string) (
				[]models.SelectOption, error) {
				return []models.SelectOption{}, nil
			},
		}
		selectOpts := &mocks.MockSelectOptionsService{
			FindAllFormSelectsFunc: func(ctx context.Context,
				language string) (*models.FormSelectsResponse, error) {
				return &models.FormSelectsResponse{}, nil
			},
		}

		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)

		columns, err := svc.BuildColumns(context.Background(), "pt")

		assert.NoError(t, err)
		assert.Len(t, columns, 15)

		cityCol := columnByLetter(t, columns, "E")
		_, ok := cityCol.FindRealValue("Rio de Janeiro")
		assert.False(t, ok)

		inNetworkCol := columnByLetter(t, columns, "H")
		real, ok := inNetworkCol.FindRealValue("Sim")
		assert.True(t, ok)
		assert.Equal(t, "true", real)
	})

	t.Run("Error - City fails", func(t *testing.T) {
		citySvc := &mocks.MockCityService{
			FindAllFunc: func(ctx context.Context, language string) (
				[]models.SelectOption, error) {
				return nil, assert.AnError
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewTemplateService(&mocks.MockCountryService{},
			citySvc, &mocks.MockSelectOptionsService{}, mockLogger)

		columns, err := svc.BuildColumns(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, columns)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Country fails", func(t *testing.T) {
		countrySvc := &mocks.MockCountryService{
			FindAllFunc: func(ctx context.Context, lang string) (
				[]models.SelectOption, error) {
				return nil, assert.AnError
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewTemplateService(countrySvc,
			&mocks.MockCityService{}, &mocks.MockSelectOptionsService{},
			mockLogger)

		columns, err := svc.BuildColumns(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, columns)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - SelectOptions fails", func(t *testing.T) {
		selectOpts := &mocks.MockSelectOptionsService{
			FindAllFormSelectsFunc: func(ctx context.Context,
				language string) (*models.FormSelectsResponse, error) {
				return nil, assert.AnError
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewTemplateService(&mocks.MockCountryService{},
			&mocks.MockCityService{}, selectOpts, mockLogger)

		columns, err := svc.BuildColumns(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, columns)
		assert.Equal(t, 1, logs.Len())
	})
}

func TestTemplateServiceCreateTemplateTable(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()

		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)

		f, err := svc.CreateTemplateTable(context.Background(), "pt")

		assert.NoError(t, err)
		if assert.NotNil(t, f) {
			assert.Equal(t, "Samples", f.GetSheetName(0))

			for cell, expected := range map[string]string{
				"A1": "origin_code",
				"E1": "city",
				"O1": "health_service_id",
			} {
				value, cellErr := f.GetCellValue("Samples", cell)
				assert.NoError(t, cellErr)
				assert.Equal(t, expected, value)
			}

			visible, visErr := f.GetSheetVisible("AcceptedValues")
			assert.NoError(t, visErr)
			assert.False(t, visible)
		}
	})

	t.Run("Error - BuildColumns fails", func(t *testing.T) {
		citySvc := &mocks.MockCityService{
			FindAllFunc: func(ctx context.Context, language string) (
				[]models.SelectOption, error) {
				return nil, assert.AnError
			},
		}
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)

		svc := services.NewTemplateService(&mocks.MockCountryService{},
			citySvc, &mocks.MockSelectOptionsService{}, mockLogger)

		f, err := svc.CreateTemplateTable(context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, f)
		assert.Equal(t, 1, logs.Len())
	})
}
