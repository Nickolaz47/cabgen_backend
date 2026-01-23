package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleAuthError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrConflictEmail):
		return http.StatusConflict, responses.RegisterEmailAlreadyExistsError
	case errors.Is(err, services.ErrConflictUsername):
		return http.StatusConflict, responses.RegisterUsernameAlreadyExistsError
	case errors.Is(err, services.ErrEmailMismatch):
		return http.StatusBadRequest, responses.RegisterValidationEmailMismatch
	case errors.Is(err, services.ErrPasswordMismatch):
		return http.StatusBadRequest, responses.RegisterValidationPasswordMismatch
	case errors.Is(err, services.ErrInvalidCountryCode):
		return http.StatusNotFound, responses.CountryNotFoundError
	case errors.Is(err, services.ErrInvalidCredentials):
		return http.StatusUnauthorized, responses.LoginInvalidCredentialsError
	case errors.Is(err, services.ErrDisabledUser):
		return http.StatusForbidden, responses.LoginInactiveUser
	case errors.Is(err, services.ErrUnauthorized):
		return http.StatusUnauthorized, responses.UnauthorizedError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
