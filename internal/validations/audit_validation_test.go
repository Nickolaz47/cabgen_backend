package validations_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/testutils"
	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestGetAuditInputFromContext(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success", func(t *testing.T) {
		expected := &models.AuditInput{
			Event:  models.AuditEventLogin,
			Source: "10.0.0.1",
		}
		c, _ := testutils.SetupGinContext(http.MethodPost, "/", "", nil, nil)
		c.Set(validations.AuditKey, expected)

		audit, ok := validations.GetAuditInputFromContext(c)

		assert.True(t, ok)
		assert.Equal(t, expected, audit)
	})

	t.Run("Error - Not In Context", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodPost, "/", "", nil, nil)

		audit, ok := validations.GetAuditInputFromContext(c)

		assert.False(t, ok)
		assert.Nil(t, audit)
	})

	t.Run("Error - Wrong Type", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodPost, "/", "", nil, nil)
		c.Set(validations.AuditKey, "not an audit")

		audit, ok := validations.GetAuditInputFromContext(c)

		assert.False(t, ok)
		assert.Nil(t, audit)
	})
}

func TestSetAuditEvent(t *testing.T) {
	testutils.SetupTestContext()

	t.Run("Success - With Metadata", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodPost, "/", "", nil, nil)
		c.Set(validations.AuditKey, &models.AuditInput{Source: "10.0.0.1"})

		validations.SetAuditEvent(c, models.AuditEventLoginFailed,
			map[string]string{"auth_identity": "<script>x</script>nick@mail.com"})

		audit, ok := validations.GetAuditInputFromContext(c)

		assert.True(t, ok)
		assert.Equal(t, models.AuditEventLoginFailed, audit.Event)
		assert.NotNil(t, audit.Metadata)

		var metadata map[string]string
		err := json.Unmarshal([]byte(*audit.Metadata), &metadata)

		assert.NoError(t, err)
		assert.False(t, strings.Contains(metadata["auth_identity"], "<script>"))
		assert.True(t, strings.Contains(metadata["auth_identity"], "nick@mail.com"))
	})

	t.Run("Success - Without Metadata", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodPost, "/", "", nil, nil)
		c.Set(validations.AuditKey, &models.AuditInput{Source: "10.0.0.1"})

		validations.SetAuditEvent(c, models.AuditEventLogin, nil)

		audit, ok := validations.GetAuditInputFromContext(c)

		assert.True(t, ok)
		assert.Equal(t, models.AuditEventLogin, audit.Event)
		assert.Nil(t, audit.Metadata)
	})

	t.Run("Success - Without Audit In Context", func(t *testing.T) {
		c, _ := testutils.SetupGinContext(http.MethodPost, "/", "", nil, nil)

		assert.NotPanics(t, func() {
			validations.SetAuditEvent(c, models.AuditEventLogin, nil)
		})
	})
}
