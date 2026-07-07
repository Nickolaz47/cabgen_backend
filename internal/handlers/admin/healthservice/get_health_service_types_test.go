package healthservice_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/healthservice"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetHealthServiceTypes(t *testing.T) {
	testutils.SetupTestContext()

	svc := &mocks.MockHealthServiceService{}
	handler := healthservice.NewAdminHealthServiceHandler(svc)

	c, w := testutils.SetupGinContext(
		http.MethodGet,
		"/api/admin/health-service/types",
		"", nil, nil,
	)
	handler.GetHealthServiceTypes(c)

	expected := testutils.ToJSON(map[string]any{
		"data": models.HealthServiceTypes,
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
