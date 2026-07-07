package analysis_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/analysis"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAnalysisTypes(t *testing.T) {
	testutils.SetupTestContext()

	svc := &mocks.MockAdminAnalysisService{}
	handler := analysis.NewAdminAnalysisHandler(svc)

	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/admin/analysis/types", "", nil, nil,
	)
	handler.GetAnalysisTypes(c)

	expected := testutils.ToJSON(
		map[string][]models.AnalysisType{
			"data": models.AnalysisTypes,
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
