package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleSampleError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.SampleNotFoundError
	case errors.Is(err, services.ErrUnauthorized):
		return http.StatusUnauthorized, responses.UnauthorizedError
	case errors.Is(err, services.ErrInvalidCountryCode):
		return http.StatusNotFound, responses.CountryNotFoundError
	case errors.Is(err, services.ErrUserNotFound):
		return http.StatusNotFound, responses.UserNotFoundError
	case errors.Is(err, services.ErrOriginNotFound):
		return http.StatusNotFound, responses.OriginNotFoundError
	case errors.Is(err, services.ErrSampleSourceNotFound):
		return http.StatusNotFound, responses.SampleSourceNotFoundError
	case errors.Is(err, services.ErrMicroorganismNotFound):
		return http.StatusNotFound, responses.MicroorganismNotFoundError
	case errors.Is(err, services.ErrSequencerNotFound):
		return http.StatusNotFound, responses.SequencerNotFoundError
	case errors.Is(err, services.ErrLaboratoryNotFound):
		return http.StatusNotFound, responses.LaboratoryNotFoundError
	case errors.Is(err, services.ErrHealthServiceNotFound):
		return http.StatusNotFound, responses.HealthServiceNotFoundError
	case errors.Is(err, services.ErrMissingFastq1):
		return http.StatusBadRequest, responses.SampleMissingFastq1
	case errors.Is(err, services.ErrMissingFastq2):
		return http.StatusBadRequest, responses.SampleMissingFastq2
	case errors.Is(err, services.ErrMissingFiles):
		return http.StatusBadRequest, responses.SampleMissingFiles
	default:
		return http.StatusInternalServerError,
			responses.GenericInternalServerError
	}
}
