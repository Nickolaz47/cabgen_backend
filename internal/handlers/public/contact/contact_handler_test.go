package contact_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/public/contact"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/data"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils/mocks"
	testmodels "github.com/CABGenOrg/cabgen_backend/internal/testutils/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicket(t *testing.T) {
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
			CreateFunc: func(ctx context.Context, input models.CreateTicketInput) (
				*models.TicketResponse, error) {
				return &mockResponse, nil
			},
		}

		handler := contact.NewTicketHandler(svc)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/contact",
			testutils.ToJSON(models.CreateTicketInput{
				Name:        "Jão",
				Email:       "jao@mail.com",
				Institution: "Fiocruz",
				Subject:     "Wrong password",
				Message:     "Cannot access my account.",
			}),
			nil, nil,
		)

		handler.CreateTicket(c)

		expected := testutils.ToJSON(map[string]any{
			"message": "Contact ticket created successfully.",
			"data":    mockResponse,
		})

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	for _, tt := range data.CreateTicketValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			svc := &mocks.MockTicketService{}
			handler := contact.NewTicketHandler(svc)
			c, w := testutils.SetupGinContext(
				http.MethodPost, "/api/contact",
				tt.Body, nil, nil,
			)

			handler.CreateTicket(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, tt.Expected, w.Body.String())
		})
	}

	t.Run("Error - Internal Server Error", func(t *testing.T) {
		svc := &mocks.MockTicketService{
			CreateFunc: func(ctx context.Context, input models.CreateTicketInput) (
				*models.TicketResponse, error) {
				return nil, services.ErrInternal
			},
		}
		handler := contact.NewTicketHandler(svc)
		c, w := testutils.SetupGinContext(
			http.MethodPost, "/api/contact",
			testutils.ToJSON(models.CreateTicketInput{
				Name:        "Jão",
				Email:       "jao@mail.com",
				Institution: "Fiocruz",
				Subject:     "Wrong password",
				Message:     "Cannot access my account.",
			}),
			nil, nil,
		)

		handler.CreateTicket(c)

		expected := testutils.ToJSON(map[string]string{
			"error": "There was a server error. Please try again.",
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})
}
