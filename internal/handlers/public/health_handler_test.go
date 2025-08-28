package public_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	testutils.SetupTestContext()
	// Simulates the http response and creates the Gin test context
	c, w := testutils.SetupGinContext(
		http.MethodGet, "/api/health", "",
		nil, nil,
	)

	public.Health(c)

	expected := `{"message": "API ok."}`

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}
