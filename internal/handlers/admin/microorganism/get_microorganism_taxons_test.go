package microorganism_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/microorganism"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetMicroorganismTaxons(t *testing.T) {
	testutils.SetupTestContext()

	svc := &mocks.MockMicroorganismService{}
	handler := microorganism.NewAdminMicroorganismHandler(svc)

	c, w := testutils.SetupGinContext(
		http.MethodGet,
		"/api/admin/microorganism/taxons",
		"",
		nil,
		nil,
	)

	handler.GetMicroorganismTaxons(c)

	expected := testutils.ToJSON(
		map[string]any{
			"data": models.Taxons,
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
