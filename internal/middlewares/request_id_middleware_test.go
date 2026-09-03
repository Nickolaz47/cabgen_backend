package middlewares_test

import (
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/logging"
	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		w, r := testutils.SetupMiddlewareContext()

		var ctxRequestID string
		r.Use(middlewares.RequestIDMiddleware())
		r.GET("/", func(c *gin.Context) {
			ctxRequestID = logging.RequestIDFromContext(c.Request.Context())
			c.Status(http.StatusOK)
		})
		testutils.DoGetRequest(r, w)

		headerID := w.Header().Get("X-Request-ID")

		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, headerID)
		_, err := uuid.Parse(headerID)
		assert.NoError(t, err, "request_id must be a valid UUID")
		assert.Equal(t, headerID, ctxRequestID,
			"request context must carry the same ID as the header")
	})

	t.Run("Unique per request", func(t *testing.T) {
		w1, r1 := testutils.SetupMiddlewareContext()
		w2, r2 := testutils.SetupMiddlewareContext()

		r1.Use(middlewares.RequestIDMiddleware())
		r2.Use(middlewares.RequestIDMiddleware())
		testutils.AddTestGetRoute(r1, http.StatusOK)
		testutils.AddTestGetRoute(r2, http.StatusOK)

		testutils.DoGetRequest(r1, w1)
		testutils.DoGetRequest(r2, w2)

		id1 := w1.Header().Get("X-Request-ID")
		id2 := w2.Header().Get("X-Request-ID")

		assert.NotEmpty(t, id1)
		assert.NotEmpty(t, id2)
		assert.NotEqual(t, id1, id2)
	})
}
