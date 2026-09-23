package middlewares_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/middlewares"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuditMiddleware(t *testing.T) {
	testutils.SetupTestContext()

	markEvent := func(c *gin.Context) {
		validations.SetAuditEvent(c, models.AuditEventLogin, nil)
		c.Status(http.StatusCreated)
	}

	t.Run("Success - Audit Written", func(t *testing.T) {
		received := make(chan *models.AuditInput, 1)
		svc := &mocks.MockAuditService{
			CreateFunc: func(ctx context.Context, input *models.AuditInput) error {
				received <- input
				return nil
			},
		}

		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuditMiddleware(svc))
		r.POST("/", markEvent)

		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		req.RemoteAddr = "192.0.2.1:1234"
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		select {
		case audit := <-received:
			assert.Equal(t, models.AuditEventLogin, audit.Event)
			assert.Equal(t, http.StatusCreated, audit.Status)
			assert.Equal(t, "192.0.2.1", audit.Source)
			assert.Nil(t, audit.UserID)
		case <-time.After(2 * time.Second):
			t.Fatal("audit write never happened")
		}
	})

	t.Run("Success - Authenticated", func(t *testing.T) {
		userID := uuid.New()
		received := make(chan *models.AuditInput, 1)
		svc := &mocks.MockAuditService{
			CreateFunc: func(ctx context.Context, input *models.AuditInput) error {
				received <- input
				return nil
			},
		}

		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuditMiddleware(svc),
			func(c *gin.Context) {
				c.Set(validations.UserTokenKey,
					&models.UserToken{ID: userID})
				c.Next()
			})
		r.POST("/", markEvent)

		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		select {
		case audit := <-received:
			assert.Equal(t, &userID, audit.UserID)
		case <-time.After(2 * time.Second):
			t.Fatal("audit write never happened")
		}
	})

	t.Run("Success - Service Error Ignored", func(t *testing.T) {
		svc := &mocks.MockAuditService{
			CreateFunc: func(ctx context.Context, input *models.AuditInput) error {
				return context.Canceled
			},
		}

		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuditMiddleware(svc))
		r.POST("/", func(c *gin.Context) {
			validations.SetAuditEvent(c, models.AuditEventLogin, nil)
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Success - Without Event", func(t *testing.T) {
		var calls int32
		svc := &mocks.MockAuditService{
			CreateFunc: func(ctx context.Context, input *models.AuditInput) error {
				atomic.AddInt32(&calls, 1)
				return nil
			},
		}

		w, r := testutils.SetupMiddlewareContext()
		testutils.AddMiddlewares(r, middlewares.AuditMiddleware(svc))
		r.POST("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodPost, "/", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, int32(0), atomic.LoadInt32(&calls))
	})
}
