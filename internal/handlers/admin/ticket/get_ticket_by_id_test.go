package ticket_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/admin/ticket"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTicketByID(t *testing.T) {
	testutils.SetupTestContext()

	admin := testmodels.NewAdminLoginUser()
	mockTicket := testmodels.NewTicket(
		uuid.NewString(),
		"Jão",
		"jão@mail.com",
		"Fiocruz",
		"Wrong password",
		"Cannot access my account.",
		&admin,
	)
	mockResponse := mockTicket.ToResponse()

	t.Run("Success", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.TicketResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := ticket.NewAdminTicketHandler(svc)
		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: mockTicket.ID.String()}})
		handler.GetTicketByID(c)

		expected := testutils.ToJSON(map[string]models.TicketResponse{
			"data": mockResponse,
		})

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Invalid ID", func(t *testing.T) {
		svc := &mocks.MockTicketService{}
		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: "abc123"}})
		handler.GetTicketByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "The URL ID is invalid.",
			},
		)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Not found", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.TicketResponse, error) {
				return nil, services.ErrNotFound
			},
		}
		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: uuid.NewString()}})
		handler.GetTicketByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "Ticket not found.",
			},
		)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Error - Internal Server", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (
				*models.TicketResponse, error) {
				return nil, services.ErrInternal
			},
		}

		handler := ticket.NewAdminTicketHandler(svc)

		c, w := testutils.SetupGinContext(http.MethodGet,
			"/api/admin/ticket", "", nil,
			gin.Params{{Key: "ticketId", Value: uuid.NewString()}})
		handler.GetTicketByID(c)

		expected := testutils.ToJSON(
			map[string]string{
				"error": "There was a server error. Please try again.",
			},
		)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
