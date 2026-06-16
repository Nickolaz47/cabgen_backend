package ticket

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminTicketHandler struct {
	Service services.TicketService
}

func NewAdminTicketHandler(svc services.TicketService) *AdminTicketHandler {
	return &AdminTicketHandler{
		Service: svc,
	}
}

func (h *AdminTicketHandler) GetTickets(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	statusFilter := c.Query("status")

	tickets, err := h.Service.FindAll(c.Request.Context(), statusFilter)
	if err != nil {
		code, errMsg := handlererrors.HandleTicketError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: tickets})
}

func (h *AdminTicketHandler) GetTicketByID(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("ticketId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	ticket, err := h.Service.FindByID(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleTicketError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: ticket})
}

func (h *AdminTicketHandler) Assign(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("ticketId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	userToken, ok := validations.GetUserTokenFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.UnauthorizedError),
		})
		return
	}

	ticket, err := h.Service.Assign(c.Request.Context(), id, userToken.ID)
	if err != nil {
		code, errMsg := handlererrors.HandleTicketError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: ticket})
}

func (h *AdminTicketHandler) Resolve(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("ticketId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	ticket, err := h.Service.Resolve(c.Request.Context(), id)
	if err != nil {
		code, errMsg := handlererrors.HandleTicketError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{Data: ticket})
}

func (h *AdminTicketHandler) DeleteTicket(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)
	rawID := c.Param("ticketId")

	id, err := uuid.Parse(rawID)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: responses.GetResponse(localizer, responses.InvalidURLID),
		})
		return
	}

	if err = h.Service.Delete(c.Request.Context(), id); err != nil {
		code, errMsg := handlererrors.HandleTicketError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusOK, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.TicketDelete),
	})
}
