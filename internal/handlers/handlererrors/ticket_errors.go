package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleTicketError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.TicketNotFoundError
	case errors.Is(err, services.ErrDeleteActiveTicket):
		return http.StatusBadRequest, responses.TicketDeleteInProgressError
	case errors.Is(err, services.ErrTicketAlreadyResolvedStatus):
		return http.StatusBadRequest, responses.TicketAlreadyResolvedError
	case errors.Is(err, services.ErrTicketIsNotOpen):
		return http.StatusBadRequest, responses.TickedIsNotOpenError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
