package contact

import (
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/translation"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	Service services.TicketService
}

func NewTicketHandler(svc services.TicketService) *TicketHandler {
	return &TicketHandler{
		Service: svc,
	}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	localizer := translation.GetLocalizerFromContext(c)

	var newTicket models.CreateTicketInput
	if errMsg, valid := validations.Validate(c, localizer, &newTicket); !valid {
		c.JSON(http.StatusBadRequest, responses.APIResponse{
			Error: errMsg,
		})
		return
	}

	ticket, err := h.Service.Create(c.Request.Context(), newTicket)
	if err != nil {
		code, errMsg := handlererrors.HandleTicketError(err)
		c.JSON(code, responses.APIResponse{
			Error: responses.GetResponse(localizer, errMsg),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.APIResponse{
		Message: responses.GetResponse(localizer, responses.TicketCreationSuccess),
		Data:    ticket,
	})
}
