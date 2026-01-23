package health_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/health"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	testutils.SetupTestContext()
	handler := health.NewHealthHandler()

	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/health", "",
		nil, nil,
	)
	handler.Health(c)

	expected := testutils.ToJSON(
		map[string]string{
			"message": "API ok.",
		},
	)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
