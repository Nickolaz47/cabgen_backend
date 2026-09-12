package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGlobalRateLimitPerMinute(t *testing.T) {
	testutils.SetupTestContext()

	doPost := func(r http.Handler) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		r.ServeHTTP(w, req)
		return w
	}

	t.Run("Within Limit", func(t *testing.T) {
		_, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.GlobalRateLimitPerMinute(5))
		r.POST("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		for range 5 {
			w := doPost(r)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})

	t.Run("Exceeds Limit", func(t *testing.T) {
		_, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.GlobalRateLimitPerMinute(1))
		r.POST("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := doPost(r)
		assert.Equal(t, http.StatusOK, w.Code)

		w = doPost(r)
		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.JSONEq(t,
			testutils.ToJSON(map[string]string{
				"error": "Too many requests. Please try again later.",
			}),
			w.Body.String())
	})

	t.Run("Independent Instances", func(t *testing.T) {
		_, r := testutils.SetupMiddlewareContext()
		r.POST("/route-a",
			middlewares.GlobalRateLimitPerMinute(1),
			func(c *gin.Context) { c.Status(http.StatusOK) })
		r.POST("/route-b",
			middlewares.GlobalRateLimitPerMinute(1),
			func(c *gin.Context) { c.Status(http.StatusOK) })

		for _, path := range []string{"/route-a", "/route-b"} {
			req, _ := http.NewRequest(http.MethodPost, path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		}
	})
}
