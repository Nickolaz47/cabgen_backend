package validations_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestIsAllowedTableFile(t *testing.T) {
	t.Run("Success - Lowercase xlsx", func(t *testing.T) {
		assert.True(t, validations.IsAllowedTableFile("amostra.xlsx"))
	})

	t.Run("Success - Uppercase xlsx", func(t *testing.T) {
		assert.True(t, validations.IsAllowedTableFile("AMOSTRA.XLSX"))
	})

	t.Run("Error - Unsupported extension", func(t *testing.T) {
		assert.False(t, validations.IsAllowedTableFile("amostra.xlsm"))
		assert.False(t, validations.IsAllowedTableFile("amostra.csv"))
	})

	t.Run("Error - No extension", func(t *testing.T) {
		assert.False(t, validations.IsAllowedTableFile("amostra"))
	})

	t.Run("Error - Name too long", func(t *testing.T) {
		name := strings.Repeat("a", 256) + ".xlsx"
		assert.False(t, validations.IsAllowedTableFile(name))
	})
}

type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func TestReadTablePart(t *testing.T) {
	t.Run("Success - Empty", func(t *testing.T) {
		data, err := validations.ReadTablePart(bytes.NewReader(nil))

		assert.NoError(t, err)
		assert.Empty(t, data)
	})

	t.Run("Success - Below limit", func(t *testing.T) {
		content := bytes.Repeat([]byte("x"), 1024)

		data, err := validations.ReadTablePart(bytes.NewReader(content))

		assert.NoError(t, err)
		assert.Equal(t, content, data)
	})

	t.Run("Success - Exactly at limit", func(t *testing.T) {
		content := bytes.Repeat([]byte("x"), int(validations.MaxTableSize))

		data, err := validations.ReadTablePart(bytes.NewReader(content))

		assert.NoError(t, err)
		assert.Len(t, data, int(validations.MaxTableSize))
	})

	t.Run("Error - Above limit", func(t *testing.T) {
		content := bytes.Repeat(
			[]byte("x"), int(validations.MaxTableSize)+1)

		data, err := validations.ReadTablePart(bytes.NewReader(content))

		assert.ErrorIs(t, err, validations.ErrTableTooLarge)
		assert.Nil(t, data)
	})

	t.Run("Error - Reader fails", func(t *testing.T) {
		data, err := validations.ReadTablePart(errReader{})

		assert.ErrorContains(t, err, "read failed")
		assert.Nil(t, data)
	})
}

var _ io.Reader = errReader{}
