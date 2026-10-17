package services_test

import (
	"context"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap/zapcore"
)

var (
	testOriginID        = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	testSampleSourceID  = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	testMicroorganismID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	testSequencerID     = uuid.MustParse("44444444-4444-4444-4444-444444444444")
	testLaboratoryID    = uuid.MustParse("55555555-5555-5555-5555-555555555555")
	testHealthServiceID = uuid.MustParse("66666666-6666-6666-6666-666666666666")
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
	originID := testOriginID
	selectOpts := &mocks.MockSelectOptionsService{
		FindAllFormSelectsFunc: func(ctx context.Context,
			language string) (*models.FormSelectsResponse, error) {
			return &models.FormSelectsResponse{
				Laboratories: []models.SelectOption{
					{Label: "LACEN/RJ", Value: testLaboratoryID.String()},
				},
				Sequencers: []models.SelectOption{
					{Label: "Illumina", Value: testSequencerID.String()},
				},
				HealthServices: []models.SelectOption{
					{Label: "Hospital Central", Value: testHealthServiceID.String()},
				},
				Origins: []models.SelectOption{
					{Label: "Humano", Value: originID.String()},
				},
				Microorganisms: []models.SelectOption{
					{Label: "Escherichia coli", Value: testMicroorganismID.String()},
				},
				SampleSources: []models.SelectOption{
					{Label: "Aspirado", Value: testSampleSourceID.String()},
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

// Labels vary by language; values stay identical across languages
func templateSuccessMocksByLang() (
	*mocks.MockCityService, *mocks.MockCountryService,
	*mocks.MockSelectOptionsService) {
	countryLabels := map[string]string{
		"pt": "Brasil", "en": "Brazil", "es": "Brasil",
	}
	genderLabels := map[string]string{
		"pt": "Masculino", "en": "Male", "es": "Masculino",
	}
	originLabels := map[string]string{
		"pt": "Humano", "en": "Human", "es": "Humano",
	}

	citySvc := &mocks.MockCityService{
		FindAllFunc: func(ctx context.Context, language string) (
			[]models.SelectOption, error) {
			return []models.SelectOption{
				{Label: "Rio de Janeiro", Value: "Rio de Janeiro"},
			}, nil
		},
	}
	countrySvc := &mocks.MockCountryService{
		FindAllFunc: func(ctx context.Context, lang string) (
			[]models.SelectOption, error) {
			return []models.SelectOption{
				{Label: countryLabels[lang], Value: "BRA"},
			}, nil
		},
	}
	selectOpts := &mocks.MockSelectOptionsService{
		FindAllFormSelectsFunc: func(ctx context.Context,
			language string) (*models.FormSelectsResponse, error) {
			return &models.FormSelectsResponse{
				Laboratories: []models.SelectOption{
					{Label: "LACEN/RJ", Value: testLaboratoryID.String()},
				},
				Sequencers: []models.SelectOption{
					{Label: "Illumina", Value: testSequencerID.String()},
				},
				HealthServices: []models.SelectOption{
					{Label: "Hospital Central", Value: testHealthServiceID.String()},
				},
				Origins: []models.SelectOption{
					{Label: originLabels[language], Value: testOriginID.String()},
				},
				Microorganisms: []models.SelectOption{
					{Label: "Escherichia coli", Value: testMicroorganismID.String()},
				},
				SampleSources: []models.SelectOption{
					{Label: "Aspirado", Value: testSampleSourceID.String()},
				},
				Genders: []models.SelectOption{
					{Label: genderLabels[language], Value: "Male"},
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

		columns, columnHeaders, err := svc.BuildColumns(
			context.Background(), "pt")

		assert.NoError(t, err)
		assert.Equal(t, "pt", gotLanguage)
		assert.Len(t, columnHeaders, 15)
		assert.Contains(t, columnHeaders, "origin_code")
		assert.Contains(t, columnHeaders, "health_service_id")
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

	t.Run("Success - Descriptions and required flags", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)

		columns, _, err := svc.BuildColumns(
			context.Background(), "pt")

		assert.NoError(t, err)

		required := map[string]bool{
			"A": true, "B": true, "C": true, "D": true, "E": true,
			"F": false, "G": false, "H": true, "I": true, "J": true,
			"K": true, "L": true, "M": true, "N": true, "O": true,
		}
		for i := range columns {
			col := columns[i]
			assert.NotEmpty(t, col.Description, "description for %s",
				col.Letter)
			assert.Equal(t, required[col.Letter], col.Required,
				"required flag for %s", col.Letter)
		}

		assert.Equal(t,
			"Código de identificação da amostra no laboratório de origem.",
			columnByLetter(t, columns, "A").Description)
		assert.Equal(t,
			"Se deseja incluir essa amostra na rede genômica. As amostras incluídas na rede são apresentadas no Dashboard mensal de vigilância.",
			columnByLetter(t, columns, "H").Description)
		assert.Equal(t, false, columnByLetter(t, columns, "F").Required)
		assert.Equal(t, false, columnByLetter(t, columns, "G").Required)

		enColumns, _, err := svc.BuildColumns(
			context.Background(), "en")
		assert.NoError(t, err)
		assert.Equal(t,
			"Sample identification code at the originating laboratory.",
			columnByLetter(t, enColumns, "A").Description)
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

		columns, columnHeaders, err := svc.BuildColumns(
			context.Background(), "pt")

		assert.NoError(t, err)
		assert.Len(t, columnHeaders, 15)
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

		columns, columnHeaders, err := svc.BuildColumns(
			context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, columns)
		assert.Empty(t, columnHeaders)
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

		columns, columnHeaders, err := svc.BuildColumns(
			context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, columns)
		assert.Empty(t, columnHeaders)
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

		columns, columnHeaders, err := svc.BuildColumns(
			context.Background(), "pt")

		assert.Error(t, err)
		assert.ErrorIs(t, err, services.ErrInternal)
		assert.Nil(t, columns)
		assert.Empty(t, columnHeaders)
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
				"A1": "origin_code *",
				"E1": "city *",
				"O1": "health_service_id *",
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

func sampleHeaders() []string {
	return []string{
		"origin_code", "collection_date", "run_number", "run_date",
		"city", "gender", "date_of_birth", "in_network",
		"country_code", "origin_id", "sample_source_id",
		"microorganism_id", "sequencer_id", "laboratory_id",
		"health_service_id",
	}
}

func fullSampleRow() []string {
	return []string{
		"ORIG-001",
		"2025-01-15",
		"42",
		"2025-01-16",
		"Rio de Janeiro",
		"Masculino",
		"",
		models.ToTranslatedInNetwork(true, "pt"),
		"Brasil",
		"Humano",
		"Aspirado",
		"Escherichia coli",
		"Illumina",
		"LACEN/RJ",
		"Hospital Central",
	}
}

func newSamplesSheet(t *testing.T, rows [][]string) *excelize.File {
	t.Helper()
	f := excelize.NewFile()
	if err := f.SetSheetName("Sheet1", "Samples"); err != nil {
		t.Fatalf("rename sheet: %v", err)
	}
	for r, row := range rows {
		for c, val := range row {
			name, err := excelize.CoordinatesToCellName(c+1, r+1)
			if err != nil {
				t.Fatalf("cell name: %v", err)
			}
			if err := f.SetCellStr("Samples", name, val); err != nil {
				t.Fatalf("set cell: %v", err)
			}
		}
	}

	return f
}

func TestTemplateServiceValidateTemplateTable(t *testing.T) {
	t.Run("Success - Full row", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		f := newSamplesSheet(t, [][]string{sampleHeaders(),
			fullSampleRow()})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		if !assert.Len(t, inputs, 1) {
			return
		}
		in := inputs[0]
		assert.Equal(t, "ORIG-001", in.OriginCode)
		assert.Equal(t, "42", in.RunNumber)
		assert.Equal(t, "2025-01-15", in.CollectionDate.Format("2006-01-02"))
		assert.Equal(t, "2025-01-16", in.RunDate.Format("2006-01-02"))
		assert.Nil(t, in.DateOfBirth)
		assert.Equal(t, "Rio de Janeiro", in.City)
		if assert.NotNil(t, in.Gender) {
			assert.Equal(t, models.Male, *in.Gender)
		}
		if assert.NotNil(t, in.InNetwork) {
			assert.True(t, *in.InNetwork)
		}
		assert.Equal(t, "BRA", in.CountryCode)
		assert.Equal(t, testOriginID, in.OriginID)
		assert.Equal(t, testSampleSourceID, in.SampleSourceID)
		assert.Equal(t, testMicroorganismID, in.MicroorganismID)
		assert.Equal(t, testSequencerID, in.SequencerID)
		assert.Equal(t, testLaboratoryID, in.LaboratoryID)
		assert.Equal(t, testHealthServiceID, in.HealthServiceID)
	})

	t.Run("Success - DD/MM dates", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		row := fullSampleRow()
		row[1] = "22/05/2019"
		row[3] = "01/06/2019"
		row[6] = "15/03/1990"
		f := newSamplesSheet(t, [][]string{sampleHeaders(), row})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		if assert.Len(t, inputs, 1) {
			in := inputs[0]
			assert.Equal(t, "2019-05-22",
				in.CollectionDate.Format("2006-01-02"))
			assert.Equal(t, "2019-06-01", in.RunDate.Format("2006-01-02"))
			if assert.NotNil(t, in.DateOfBirth) {
				assert.Equal(t, "1990-03-15",
					in.DateOfBirth.Format("2006-01-02"))
			}
		}
	})

	t.Run("Success - Date serial", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		f := newSamplesSheet(t, [][]string{sampleHeaders(),
			fullSampleRow()})
		if err := f.SetCellValue("Samples", "B2", 43607); err != nil {
			t.Fatalf("set serial: %v", err)
		}

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		if assert.Len(t, inputs, 1) {
			assert.Equal(t, "2019-05-22",
				inputs[0].CollectionDate.Format("2006-01-02"))
		}
	})

	t.Run("Success - Ambiguous text reads as DD/MM", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		row := fullSampleRow()
		row[1] = "05/06/2019"
		f := newSamplesSheet(t, [][]string{sampleHeaders(), row})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		if assert.Len(t, inputs, 1) {
			assert.Equal(t, "2019-06-05",
				inputs[0].CollectionDate.Format("2006-01-02"))
		}
	})

	t.Run("Error - Invalid date", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		row := fullSampleRow()
		row[1] = "abc"
		f := newSamplesSheet(t, [][]string{sampleHeaders(), row})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		var tvErr *services.TableValueError
		if assert.ErrorAs(t, err, &tvErr) {
			assert.Equal(t, 2, tvErr.Row)
			assert.Equal(t, "collection_date", tvErr.Col)
		}
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Empty row skipped", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		emptyRow := make([]string, len(sampleHeaders()))
		f := newSamplesSheet(t, [][]string{sampleHeaders(), emptyRow,
			fullSampleRow()})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		assert.Len(t, inputs, 1)
	})

	t.Run("Success - Header case and spaces tolerated", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		headers := make([]string, 0, len(sampleHeaders()))
		for _, h := range sampleHeaders() {
			headers = append(headers, "  "+strings.ToUpper(h)+"  ")
		}
		f := newSamplesSheet(t, [][]string{headers, fullSampleRow()})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		assert.Len(t, inputs, 1)
	})

	t.Run("Success - Header with required asterisk accepted", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		headers := make([]string, 0, len(sampleHeaders()))
		for _, h := range sampleHeaders() {
			headers = append(headers, h+" *")
		}
		f := newSamplesSheet(t, [][]string{headers, fullSampleRow()})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		assert.Len(t, inputs, 1)
	})

	t.Run("Success - Short row", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		f := newSamplesSheet(t, [][]string{sampleHeaders(),
			{"ORIG-001"}})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		if assert.Len(t, inputs, 1) {
			assert.Equal(t, "ORIG-001", inputs[0].OriginCode)
			assert.Empty(t, inputs[0].City)
			assert.Nil(t, inputs[0].InNetwork)
		}
	})

	t.Run("Error - Sheet missing", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		f := excelize.NewFile()

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.ErrorIs(t, err, services.ErrInvalidTable)
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Empty sheet", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		f := newSamplesSheet(t, nil)

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.ErrorIs(t, err, services.ErrEmptyTable)
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Missing header", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		headers := make([]string, 0, len(sampleHeaders()))
		for _, h := range sampleHeaders() {
			if h != "city" {
				headers = append(headers, h)
			}
		}
		f := newSamplesSheet(t, [][]string{headers, fullSampleRow()})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		var thErr *services.TableHeadersError
		if assert.ErrorAs(t, err, &thErr) {
			assert.Equal(t, []string{"city"}, thErr.Invalid)
		}
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Unknown header", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		headers := append(sampleHeaders(), "foo")
		f := newSamplesSheet(t, [][]string{headers, fullSampleRow()})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		var thErr *services.TableHeadersError
		if assert.ErrorAs(t, err, &thErr) {
			assert.Equal(t, []string{"foo"}, thErr.Invalid)
		}
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Dropdown value outside list", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		row := fullSampleRow()
		row[4] = "Nowhere"
		f := newSamplesSheet(t, [][]string{sampleHeaders(), row})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		var tvErr *services.TableValueError
		if assert.ErrorAs(t, err, &tvErr) {
			assert.Equal(t, 2, tvErr.Row)
			assert.Equal(t, "city", tvErr.Col)
		}
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Error - Too many rows", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)

		rows := [][]string{sampleHeaders()}
		for i := 0; i < 1001; i++ {
			rows = append(rows, fullSampleRow())
		}
		f := newSamplesSheet(t, rows)

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.ErrorIs(t, err, services.ErrTooManyRows)
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})

	t.Run("Success - Max rows boundary", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocks()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)

		rows := [][]string{sampleHeaders()}
		for i := 0; i < 1000; i++ {
			rows = append(rows, fullSampleRow())
		}
		f := newSamplesSheet(t, rows)

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		assert.Len(t, inputs, 1000)
	})

	t.Run("Success - Labels in another language", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocksByLang()
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			nil)
		row := fullSampleRow()
		row[5] = "Male"
		row[7] = models.ToTranslatedInNetwork(true, "en")
		row[8] = "Brazil"
		row[9] = "Human"
		f := newSamplesSheet(t, [][]string{sampleHeaders(), row})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		assert.NoError(t, err)
		if assert.Len(t, inputs, 1) {
			in := inputs[0]
			assert.Equal(t, "BRA", in.CountryCode)
			if assert.NotNil(t, in.InNetwork) {
				assert.True(t, *in.InNetwork)
			}
			if assert.NotNil(t, in.Gender) {
				assert.Equal(t, models.Male, *in.Gender)
			}
			assert.Equal(t, testOriginID, in.OriginID)
		}
	})

	t.Run("Error - Unknown in every language", func(t *testing.T) {
		citySvc, countrySvc, selectOpts := templateSuccessMocksByLang()
		mockLogger, logs := testutils.NewMockLogger(zapcore.ErrorLevel)
		svc := services.NewTemplateService(countrySvc, citySvc, selectOpts,
			mockLogger)
		row := fullSampleRow()
		row[8] = "Atlantis"
		f := newSamplesSheet(t, [][]string{sampleHeaders(), row})

		inputs, err := svc.ValidateTemplateTable(context.Background(), "pt",
			f)

		var tvErr *services.TableValueError
		if assert.ErrorAs(t, err, &tvErr) {
			assert.Equal(t, 2, tvErr.Row)
			assert.Equal(t, "country_code", tvErr.Col)
		}
		assert.Nil(t, inputs)
		assert.Equal(t, 1, logs.Len())
	})
}
