package sample_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/common/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func templateFile(t *testing.T) *excelize.File {
	t.Helper()

	f := excelize.NewFile()
	assert.NoError(t, f.SetSheetName("Sheet1", "Samples"))
	assert.NoError(t, f.SetCellValue("Samples", "A1", "city"))

	return f
}

func TestDownloadSampleTemplate(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			CreateTemplateTableFunc: func(ctx context.Context,
				language string) (*excelize.File, error) {
				assert.Equal(t, "en", language)
				return templateFile(t), nil
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{},
			tmpl)

		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/samples/template", "", nil, nil)

		handler.DownloadSampleTemplate(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t,
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			w.Header().Get("Content-Type"))
		assert.Equal(t,
			"attachment; filename=cabgen_samples_template.xlsx",
			w.Header().Get("Content-Disposition"))

		f, err := excelize.OpenReader(bytes.NewReader(w.Body.Bytes()))
		assert.NoError(t, err)
		defer f.Close()

		rows, err := f.GetRows("Samples")
		assert.NoError(t, err)
		assert.Equal(t, [][]string{{"city"}}, rows)
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		tmpl := &mocks.MockTemplateService{
			CreateTemplateTableFunc: func(ctx context.Context,
				language string) (*excelize.File, error) {
				return nil, services.ErrInternal
			},
		}
		handler := sample.NewSampleHandler(&mocks.MockSampleService{},
			tmpl)

		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/samples/template", "", nil, nil)

		handler.DownloadSampleTemplate(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t,
			`{"error":"There was a server error. Please try again."}`,
			w.Body.String())
	})
}
