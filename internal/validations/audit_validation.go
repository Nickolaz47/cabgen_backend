package validations

import (
	"encoding/json"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/gin-gonic/gin"
	sanitize "github.com/mrz1836/go-sanitize"
)

const AuditKey = "Audit"

func GetAuditInputFromContext(c *gin.Context) (*models.AuditInput, bool) {
	rawAudit, exists := c.Get(AuditKey)
	if !exists {
		return nil, false
	}

	audit, ok := rawAudit.(*models.AuditInput)
	if !ok {
		return nil, false
	}

	return audit, true
}

func SetAuditEvent(c *gin.Context, event string, metadata map[string]string) {
	audit, ok := GetAuditInputFromContext(c)
	if !ok {
		return
	}

	audit.Event = event

	if metadata == nil {
		audit.Metadata = nil
		return
	}

	sanitized := make(map[string]string, len(metadata))
	for key, value := range metadata {
		sanitized[key] = sanitize.XSS(value)
	}

	if raw, err := json.Marshal(sanitized); err == nil {
		metadataStr := string(raw)
		audit.Metadata = &metadataStr
	}
}
