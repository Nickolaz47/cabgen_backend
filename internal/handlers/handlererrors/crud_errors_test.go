package handlererrors_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHandleCrudError(t *testing.T) {
	testutils.SetupTestContext()

	handle := handlererrors.HandleCrudError(handlererrors.CrudErrorMessages{
		Conflict: responses.CountryAlreadyExistsError,
		NotFound: responses.CountryNotFoundError,
	})

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{"Conflict", services.ErrConflict, http.StatusConflict,
			responses.CountryAlreadyExistsError},
		{"NotFound", services.ErrNotFound, http.StatusNotFound,
			responses.CountryNotFoundError},
		{"Default", errors.New("unknown"), http.StatusInternalServerError,
			responses.GenericInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := handle(tt.err)
			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantMsg, msg)
		})
	}
}

func TestDomainCrudErrors(t *testing.T) {
	testutils.SetupTestContext()

	handlers := []struct {
		name        string
		handle      func(err error) (int, string)
		conflictMsg string
		notFoundMsg string
	}{
		{"Country", handlererrors.HandleCountryError,
			responses.CountryAlreadyExistsError,
			responses.CountryNotFoundError},
		{"Laboratory", handlererrors.HandleLaboratoryError,
			responses.LaboratoryNameAlreadyExistsError,
			responses.LaboratoryNotFoundError},
		{"Sequencer", handlererrors.HandleSequencerError,
			responses.SequencerModelAlreadyExistsError,
			responses.SequencerNotFoundError},
		{"Origin", handlererrors.HandleOriginError,
			responses.OriginAlreadyExistsError,
			responses.OriginNotFoundError},
		{"SampleSource", handlererrors.HandleSampleSourceError,
			responses.SampleSourceAlreadyExistsError,
			responses.SampleSourceNotFoundError},
		{"Microorganism", handlererrors.HandleMicroorganismError,
			responses.MicroorganismAlreadyExistsError,
			responses.MicroorganismNotFoundError},
		{"HealthService", handlererrors.HandleHealthServiceError,
			responses.HealthServiceNameAlreadyExistsError,
			responses.HealthServiceNotFoundError},
	}

	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			status, msg := h.handle(services.ErrConflict)
			assert.Equal(t, http.StatusConflict, status)
			assert.Equal(t, h.conflictMsg, msg)

			status, msg = h.handle(services.ErrNotFound)
			assert.Equal(t, http.StatusNotFound, status)
			assert.Equal(t, h.notFoundMsg, msg)

			status, msg = h.handle(errors.New("unknown"))
			assert.Equal(t, http.StatusInternalServerError, status)
			assert.Equal(t, responses.GenericInternalServerError, msg)
		})
	}
}
