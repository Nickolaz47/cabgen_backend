package handlererrors_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/handlers/handlererrors"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleAnalysisError(t *testing.T) {
	testutils.SetupTestContext()

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"NotFound", services.ErrNotFound, http.StatusNotFound},
		{"ExceededDownloadLimit", services.ErrExceededDownloadLimit, http.StatusBadRequest},
		{"FastQCDownload", services.ErrFastQCDownload, http.StatusBadRequest},
		{"ZipNotFound", services.ErrZipNotFound, http.StatusNotFound},
		{"MissingFiles", services.ErrMissingFiles, http.StatusBadRequest},
		{"Unauthorized", services.ErrUnauthorized, http.StatusUnauthorized},
		{"SampleNotFound", services.ErrSampleNotFound, http.StatusNotFound},
		{"UserNotFound", services.ErrUserNotFound, http.StatusNotFound},
		{"MissingFastq1", services.ErrMissingFastq1, http.StatusBadRequest},
		{"MissingFastq2", services.ErrMissingFastq2, http.StatusBadRequest},
		{"DeleteRunningAnalysis", services.ErrDeleteRunningAnalysis, http.StatusBadRequest},
		{"Default", errors.New("unknown"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := handlererrors.HandleAnalysisError(tt.err)
			assert.Equal(t, tt.wantStatus, status)
			assert.NotEmpty(t, msg)
		})
	}
}

func TestRespondHTMLError(t *testing.T) {
	testutils.SetupTestContext()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlererrors.RespondHTMLError(c, http.StatusNotFound, "Not found")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
	assert.True(t, strings.Contains(w.Body.String(), "Not found"))
}
