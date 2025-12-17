package sequencer

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

func handleError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrConflict):
		return http.StatusConflict, responses.SequencerModelAlreadyExistsError
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.SequencerNotFoundError
	default:
		return http.StatusInternalServerError, responses.GenericInternalServerError
	}
}
