package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleAuditError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrInternal):
		return http.StatusInternalServerError, responses.GenericInternalServerError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
