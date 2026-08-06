package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleUserError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrConflictUsername):
		return http.StatusConflict, responses.RegisterUsernameAlreadyExistsError
	case errors.Is(err, services.ErrConflictEmail):
		return http.StatusConflict, responses.RegisterEmailAlreadyExistsError
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.UserNotFoundError
	case errors.Is(err, services.ErrInvalidCountryCode):
		return http.StatusNotFound, responses.CountryNotFoundError
	case errors.Is(err, services.ErrCurrentPasswordMismatch):
		return http.StatusBadRequest, responses.CurrentPasswordMismatchError
	case errors.Is(err, services.ErrInvalidEmailUpdateToken):
		return http.StatusBadRequest, responses.InvalidEmailUpdateTokenError
	case errors.Is(err, services.ErrExpiredEmailUpdateToken):
		return http.StatusBadRequest, responses.EmailUpdateTokenExpiredError
	case errors.Is(err, services.ErrEmailSame):
		return http.StatusBadRequest, responses.EmailSameError
	case errors.Is(err, services.ErrDuplicateTask):
		return http.StatusConflict, responses.DuplicateTaskError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
