package sample_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/sample"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetSampleGenders(t *testing.T) {
	testutils.SetupTestContext()

	svc := &mocks.MockSampleService{}
	handler := sample.NewAdminSampleHandler(svc)

	c, w := testutils.SetupGinContext(
		http.MethodGet,
		"/api/admin/sample/genders",
		"",
		nil,
		nil,
	)

	handler.GetSampleGenders(c)

	expected := testutils.ToJSON(
		map[string]any{
			"data": models.Genders,
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
