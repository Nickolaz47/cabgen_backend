package utils_test

import (
	"testing"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

type mockItem struct {
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	Internal  string    `json:"-"`
	NoTag     string
}

type mockItemWithPointer struct {
	Name  string   `json:"name"`
	Score *float64 `json:"score"`
}

type mockItemWithJSON struct {
	Name string `json:"name"`
	Meta []byte `json:"meta"`
}

func ptr[T any](v T) *T {
	return &v
}

func TestGenerateDynamicTSV(t *testing.T) {
	fixedTime := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)

	t.Run("Success - Single item", func(t *testing.T) {
		items := []mockItem{
			{Name: "Alice", Age: 30, CreatedAt: fixedTime, Internal: "skip", NoTag: "skip"},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), "name\tage\tcreated_at")
		assert.Contains(t, string(result), "Alice\t30\t15-06-2024 10:30:00")
		assert.NotContains(t, string(result), "skip")
	})

	t.Run("Success - Multiple items", func(t *testing.T) {
		items := []mockItem{
			{Name: "Alice", Age: 30, CreatedAt: fixedTime},
			{Name: "Bob", Age: 25, CreatedAt: fixedTime},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), "Alice")
		assert.Contains(t, string(result), "Bob")
	})

	t.Run("Success - Pointer slice", func(t *testing.T) {
		items := []*mockItem{
			{Name: "Carol", Age: 22, CreatedAt: fixedTime},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), "Carol")
	})

	t.Run("Success - Nil pointer field", func(t *testing.T) {
		items := []mockItemWithPointer{
			{Name: "Dave", Score: nil},
			{Name: "Eve", Score: ptr(9.5)},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), "Dave\t")
		assert.Contains(t, string(result), "Eve\t9.5")
	})

	t.Run("Success - JSON bytes field", func(t *testing.T) {
		items := []mockItemWithJSON{
			{Name: "Frank", Meta: []byte(`{"key":"value"}`)},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), `"{""key"":""value""}"`)
	})
	t.Run("Success - JSON bytes null", func(t *testing.T) {
		items := []mockItemWithJSON{
			{Name: "Grace", Meta: []byte("null")},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), "Grace\t")
	})

	t.Run("Success - JSON bytes empty", func(t *testing.T) {
		items := []mockItemWithJSON{
			{Name: "Heidi", Meta: []byte{}},
		}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Contains(t, string(result), "Heidi\t")
	})

	t.Run("Success - Empty slice", func(t *testing.T) {
		items := []mockItem{}

		result, err := utils.GenerateDynamicTSV(items)

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Error - Not a slice", func(t *testing.T) {
		_, err := utils.GenerateDynamicTSV("not a slice")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "expected a slice")
	})
}
