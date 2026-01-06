package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleLaboratoryError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrConflict):
		return http.StatusConflict, responses.LaboratoryNameAlreadyExistsError
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.LaboratoryNotFoundError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
