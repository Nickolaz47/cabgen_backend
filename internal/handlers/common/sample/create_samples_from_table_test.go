package sample_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func tableXlsxBytes(t *testing.T) []byte {
	t.Helper()

	f := excelize.NewFile()
	defer f.Close()

	buf, err := f.WriteToBuffer()
	assert.NoError(t, err)

	return buf.Bytes()
}

func tableMultipart(t *testing.T, filename string,
	content []byte) (*bytes.Buffer, string) {
	t.Helper()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	fw, _ := mw.CreateFormFile("table", filename)
	_, _ = fw.Write(content)
	mw.Close()

	return &buf, mw.FormDataContentType()
}

func tableMultipartWithUser(t *testing.T, filename string,
	content []byte, userID string) (*bytes.Buffer, string) {
	t.Helper()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	if userID != "" {
		assert.NoError(t, mw.WriteField("user_id", userID))
	}
	fw, _ := mw.CreateFormFile("table", filename)
	_, _ = fw.Write(content)
	mw.Close()

	return &buf, mw.FormDataContentType()
}

func tableInputs() []models.SampleCreateInput {
	inNetwork := true
	id := uuid.New()

	input := models.SampleCreateInput{
		OriginCode:      "BR-RJ-01",
		CollectionDate:  models.Date{Time: time.Date(2026, 5, 20, 0, 0, 0, 0, time.UTC)},
		RunNumber:       "RUN-2026-XYZ",
		RunDate:         models.Date{Time: time.Date(2026, 5, 25, 0, 0, 0, 0, time.UTC)},
		City:            "Marica",
		InNetwork:       &inNetwork,
		CountryCode:     "BRA",
		OriginID:        id,
		SampleSourceID:  id,
		MicroorganismID: id,
		SequencerID:     id,
		LaboratoryID:    id,
		HealthServiceID: id,
	}

	return []models.SampleCreateInput{input, input}
}

func TestCreateSamplesFromTable(t *testing.T) {
	testutils.SetupTestContext()
	mockUserID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateManyFunc: func(ctx context.Context,
				inputs []models.SampleCreateDTO) (int, error) {
				assert.Len(t, inputs, 2)
				for _, dto := range inputs {
					assert.Equal(t, mockUserID, dto.UserID)
				}
				return 2, nil
			},
		}
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				assert.Equal(t, "en", language)
				return tableInputs(), nil
			},
		}
		handler := sample.NewSampleHandler(svc, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t,
			`{"message":"2 samples created successfully."}`,
			w.Body.String())
	})

	t.Run("Error - Invalid Headers", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return nil, &services.TableHeadersError{
					Invalid: []string{"city", "foo"}}
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Missing or invalid headers: city, foo."}`,
			w.Body.String())
	})

	t.Run("Error - Invalid Value", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return nil, &services.TableValueError{
					Row: 2, Col: "city"}
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Invalid value in row 2, column city."}`,
			w.Body.String())
	})

	t.Run("Error - Empty Table", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return nil, services.ErrEmptyTable
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"The spreadsheet contains no data rows."}`,
			w.Body.String())
	})

	t.Run("Error - Invalid Table", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return nil, services.ErrInvalidTable
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Could not read the uploaded table. Check the file and try again."}`,
			w.Body.String())
	})

	t.Run("Error - Internal Validation", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return nil, errors.New("boom")
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t,
			`{"error":"There was a server error. Please try again."}`,
			w.Body.String())
	})

	t.Run("Error - Binding", func(t *testing.T) {
		inputs := tableInputs()
		inputs[0].OriginCode = "AB"
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return inputs, nil
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]string
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.NotEmpty(t, resp["error"])
	})

	t.Run("Error - Create Many", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateManyFunc: func(ctx context.Context,
				inputs []models.SampleCreateDTO) (int, error) {
				return 0, services.ErrSampleSourceNotFound
			},
		}
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return tableInputs(), nil
			},
		}
		handler := sample.NewSampleHandler(svc, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"Sample source not found."}`,
			w.Body.String())
	})

	t.Run("Error - No Rows", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return nil, nil
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{}, tmpl)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"The spreadsheet contains no data rows."}`,
			w.Body.String())
	})

	t.Run("Error - Missing File", func(t *testing.T) {
		handler := sample.NewSampleHandler(&mocks.MockSampleService{},
			nil)

		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		fw, _ := mw.CreateFormFile("other", "notes.txt")
		_, _ = fw.Write([]byte("x"))
		mw.Close()

		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", &buf, mw.FormDataContentType(), nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error":"No table was sent."}`,
			w.Body.String())
	})

	t.Run("Error - Unsupported File", func(t *testing.T) {
		handler := sample.NewSampleHandler(&mocks.MockSampleService{},
			nil)

		buf, ctype := tableMultipart(t, "samples.csv",
			[]byte("city;country\nMarica;BRA"))

		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Unsupported table format. Send an .xlsx file."}`,
			w.Body.String())
	})

	t.Run("Error - Corrupt File", func(t *testing.T) {
		handler := sample.NewSampleHandler(&mocks.MockSampleService{},
			nil)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			[]byte("not a workbook"))

		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey, &models.UserToken{ID: mockUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Could not read the uploaded table. Check the file and try again."}`,
			w.Body.String())
	})

	t.Run("Error - Unauthorized", func(t *testing.T) {
		handler := sample.NewSampleHandler(&mocks.MockSampleService{},
			nil)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t,
			`{"error":"Unauthorized. Please log in to continue."}`,
			w.Body.String())
	})
}

func TestCreateSamplesFromTableAdminScope(t *testing.T) {
	testutils.SetupTestContext()
	tokenUserID := uuid.New()
	targetUserID := uuid.New()

	t.Run("Success - Admin User ID", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateManyFunc: func(ctx context.Context,
				inputs []models.SampleCreateDTO) (int, error) {
				assert.Len(t, inputs, 2)
				for _, dto := range inputs {
					assert.Equal(t, targetUserID, dto.UserID)
				}
				return 2, nil
			},
		}
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return tableInputs(), nil
			},
		}
		handler := sample.NewAdminSampleHandler(svc, tmpl)

		buf, ctype := tableMultipartWithUser(t, "samples.xlsx",
			tableXlsxBytes(t), targetUserID.String())
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/admin/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey,
			&models.UserToken{ID: tokenUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t,
			`{"message":"2 samples created successfully."}`,
			w.Body.String())
	})

	t.Run("Error - Missing User ID", func(t *testing.T) {
		handler := sample.NewAdminSampleHandler(
			&mocks.MockSampleService{}, nil)

		buf, ctype := tableMultipart(t, "samples.xlsx",
			tableXlsxBytes(t))
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/admin/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey,
			&models.UserToken{ID: tokenUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Invalid input. Please check the fields and try again."}`,
			w.Body.String())
	})

	t.Run("Error - Invalid User ID", func(t *testing.T) {
		handler := sample.NewAdminSampleHandler(
			&mocks.MockSampleService{}, nil)

		buf, ctype := tableMultipartWithUser(t, "samples.xlsx",
			tableXlsxBytes(t), "not-a-uuid")
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/admin/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey,
			&models.UserToken{ID: tokenUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t,
			`{"error":"Invalid input. Please check the fields and try again."}`,
			w.Body.String())
	})

	t.Run("Success - Common Ignores User ID", func(t *testing.T) {
		svc := &mocks.MockSampleService{
			CreateManyFunc: func(ctx context.Context,
				inputs []models.SampleCreateDTO) (int, error) {
				for _, dto := range inputs {
					assert.Equal(t, tokenUserID, dto.UserID)
				}
				return len(inputs), nil
			},
		}
		tmpl := &mocks.MockTemplateService{
			ValidateTemplateTableFunc: func(ctx context.Context,
				language string, file *excelize.File) (
				[]models.SampleCreateInput, error) {
				return tableInputs(), nil
			},
		}
		handler := sample.NewSampleHandler(svc, tmpl)

		buf, ctype := tableMultipartWithUser(t, "samples.xlsx",
			tableXlsxBytes(t), targetUserID.String())
		c, w := testutils.SetupGinMultipartContext(http.MethodPost,
			"/api/samples/table", buf, ctype, nil, nil)
		c.Set(validations.UserTokenKey,
			&models.UserToken{ID: tokenUserID})

		handler.CreateSamplesFromTable(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
