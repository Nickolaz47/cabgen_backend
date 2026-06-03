package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func HandleAnalysisError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.AnalysisNotFoundError
	case errors.Is(err, services.ErrExceededDownloadLimit):
		return http.StatusBadRequest, responses.AnalysisExceededLimitError
	case errors.Is(err, services.ErrUnauthorized):
		return http.StatusUnauthorized, responses.UnauthorizedError
	case errors.Is(err, services.ErrSampleNotFound):
		return http.StatusNotFound, responses.SampleNotFoundError
	case errors.Is(err, services.ErrUserNotFound):
		return http.StatusNotFound, responses.UserNotFoundError
	case errors.Is(err, services.ErrMissingFastq1):
		return http.StatusBadRequest, responses.SampleMissingFastq1
	case errors.Is(err, services.ErrMissingFastq2):
		return http.StatusBadRequest, responses.SampleMissingFastq2
	case errors.Is(err, services.ErrDeleteRunningAnalysis):
		return http.StatusBadRequest, responses.AnalysisDeleteRunningError
	default:
		return http.StatusInternalServerError,
			responses.GenericInternalServerError
	}
}
