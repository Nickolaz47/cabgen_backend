package ticket_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/ticket"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteTicket(t *testing.T) {
	testutils.SetupTestContext()

	ticketID := uuid.NewString()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return nil
			},
		}

		handler := ticket.NewAdminTicketHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodDelete,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: ticketID}})

		handler.DeleteTicket(c)

		expected := testutils.ToJSON(map[string]string{
			"message": "Ticket deleted successfully.",
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockTicketService{}
		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodDelete,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: "abc123"}})

		handler.DeleteTicket(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Ticket In Progress", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return services.ErrDeleteActiveTicket
			},
		}
		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodDelete,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: ticketID}})

		handler.DeleteTicket(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Cannot delete a ticket that is in progress.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not Found", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return services.ErrNotFound
			},
		}
		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodDelete,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: ticketID}})

		handler.DeleteTicket(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Ticket not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server Error", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			DeleteFunc: func(ctx context.Context, id uuid.UUID) error {
				return services.ErrInternal
			},
		}
		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodDelete,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: ticketID}})

		handler.DeleteTicket(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
