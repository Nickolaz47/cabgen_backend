package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleMicroorganismError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrConflict):
		return http.StatusConflict, responses.MicroorganismAlreadyExistsError
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.MicroorganismNotFoundError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
