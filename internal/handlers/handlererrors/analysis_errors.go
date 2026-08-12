package handlererrors

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/CABGenOrg/cabgen_backend/internal/responses"
	"github.com/CABGenOrg/cabgen_backend/internal/services"
	"github.com/gin-gonic/gin"
)

func HandleAnalysisError(err error) (int, string) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		return http.StatusNotFound, responses.AnalysisNotFoundError
	case errors.Is(err, services.ErrExceededDownloadLimit):
		return http.StatusBadRequest, responses.AnalysisExceededLimitError
	case errors.Is(err, services.ErrFastQCDownload):
		return http.StatusBadRequest, responses.AnalysisFastQCDownloadError
	case errors.Is(err, services.ErrZipNotFound):
		return http.StatusNotFound, responses.AnalysisZipNotFound
	case errors.Is(err, services.ErrUnauthorized):
		return http.StatusUnauthorized, responses.UnauthorizedError
	case errors.Is(err, services.ErrSampleNotFound):
		return http.StatusNotFound, responses.SampleNotFoundError
	case errors.Is(err, services.ErrUserNotFound):
		return http.StatusNotFound, responses.UserNotFoundError
	case errors.Is(err, services.ErrMissingFiles):
		return http.StatusBadRequest, responses.SampleMissingFiles
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

var errorPageTemplate = template.Must(template.New(
	"error").Parse(`<!DOCTYPE html>
<html lang="pt-br">
<head>
	<meta charset="utf-8">
	<title>HTTP {{.Code}}</title>
	<style>
		body { font-family: sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #f5f5f5; }
		.box { text-align: center; padding: 2rem; }
		h1 { font-size: 3rem; margin-bottom: 0.5rem; color: #333; }
		p { color: #666; }
	</style>
</head>
<body>
	<div class="box">
		<h1>{{.Code}}</h1>
		<p>{{.Message}}</p>
	</div>
</body>
</html>`))

type errorPageData struct {
	Code    int
	Message string
}

func RespondHTMLError(c *gin.Context, code int, message string) {
	c.Status(code)
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = errorPageTemplate.Execute(c.Writer, errorPageData{Code: code,
		Message: message})
}
