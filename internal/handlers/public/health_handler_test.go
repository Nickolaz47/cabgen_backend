package public_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	testutils.SetupTestContext()
	// Simulates the http response and creates the Gin test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	public.Health(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "API ok."}`, w.Body.String())
}
