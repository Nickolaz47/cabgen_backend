package handlererrors

import (
	"errors"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
)

type CrudErrorMessages struct {
	Conflict string
	NotFound string
}

func HandleCrudError(msgs CrudErrorMessages) func(err error) (int, string) {
	return func(err error) (int, string) {
		switch {
		case errors.Is(err, services.ErrConflict):
			return http.StatusConflict, msgs.Conflict
		case errors.Is(err, services.ErrNotFound):
			return http.StatusNotFound, msgs.NotFound
		default:
			return http.StatusInternalServerError,
				responses.GenericInternalServerError
		}
	}
}

var HandleCountryError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.CountryAlreadyExistsError,
	NotFound: responses.CountryNotFoundError,
})

var HandleLaboratoryError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.LaboratoryNameAlreadyExistsError,
	NotFound: responses.LaboratoryNotFoundError,
})

var HandleSequencerError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.SequencerModelAlreadyExistsError,
	NotFound: responses.SequencerNotFoundError,
})

var HandleOriginError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.OriginAlreadyExistsError,
	NotFound: responses.OriginNotFoundError,
})

var HandleSampleSourceError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.SampleSourceAlreadyExistsError,
	NotFound: responses.SampleSourceNotFoundError,
})

var HandleMicroorganismError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.MicroorganismAlreadyExistsError,
	NotFound: responses.MicroorganismNotFoundError,
})

var HandleHealthServiceError = HandleCrudError(CrudErrorMessages{
	Conflict: responses.HealthServiceNameAlreadyExistsError,
	NotFound: responses.HealthServiceNotFoundError,
})