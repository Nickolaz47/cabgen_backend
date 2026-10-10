package handlererrors_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestHandleTableError(t *testing.T) {
	testutils.SetupTestContext()

	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
		wantData   map[string]any
	}{
		{"TableHeaders", &services.TableHeadersError{
			Invalid: []string{"city", "foo"}}, http.StatusBadRequest,
			responses.SampleTableHeadersError,
			map[string]any{"Param": "city, foo"}},
		{"TableValue", &services.TableValueError{Row: 2, Col: "city"},
			http.StatusBadRequest, responses.SampleTableValueError,
			map[string]any{"Row": 2, "Col": "city"}},
		{"WrappedTableHeaders",
			fmt.Errorf("validate: %w",
				&services.TableHeadersError{Invalid: []string{"city"}}),
			http.StatusBadRequest, responses.SampleTableHeadersError,
			map[string]any{"Param": "city"}},
		{"EmptyTable", services.ErrEmptyTable, http.StatusBadRequest,
			responses.SampleTableEmptyError, nil},
		{"InvalidTable", services.ErrInvalidTable, http.StatusBadRequest,
			responses.SampleTableReadError, nil},
		{"Default", errors.New("unknown"),
			http.StatusInternalServerError,
			responses.GenericInternalServerError, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg, data := handlererrors.HandleTableError(tt.err)

			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantMsg, msg)
			assert.Equal(t, tt.wantData, data)
		})
	}
}
