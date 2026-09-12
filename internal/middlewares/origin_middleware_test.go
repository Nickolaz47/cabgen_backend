package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOriginCheckMiddleware(t *testing.T) {
	testutils.SetupTestContext()

	original := config.FrontendURL
	config.FrontendURL = "http://localhost:5173"
	defer func() { config.FrontendURL = original }()

	setup := func(method, origin string) *httptest.ResponseRecorder {
		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.OriginCheckMiddleware())
		r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })
		r.POST("/", func(c *gin.Context) { c.Status(http.StatusOK) })

		req, _ := http.NewRequest(method, "/", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("GET With Wrong Origin Passes", func(t *testing.T) {
		w := setup(http.MethodGet, "http://evil.com")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET Without Origin Passes", func(t *testing.T) {
		w := setup(http.MethodGet, "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST Without Origin Passes", func(t *testing.T) {
		w := setup(http.MethodPost, "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST With Wrong Origin Rejected", func(t *testing.T) {
		w := setup(http.MethodPost, "http://evil.com")
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.JSONEq(t,
			testutils.ToJSON(map[string]string{
				"error": "Request origin not allowed.",
			}),
			w.Body.String())
	})

	t.Run("POST With Frontend Origin Passes", func(t *testing.T) {
		w := setup(http.MethodPost, "http://localhost:5173")
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
