package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTime(t *testing.T) {
	t.Run("Success - Zero", func(t *testing.T) {
		assert.Equal(t, "00:00:00", FormatTime(0))
	})

	t.Run("Success - Seconds Only", func(t *testing.T) {
		assert.Equal(t, "00:00:45", FormatTime(45))
	})

	t.Run("Success - Minutes and Seconds", func(t *testing.T) {
		assert.Equal(t, "00:05:30", FormatTime(330))
	})

	t.Run("Success - Hours Minutes Seconds", func(t *testing.T) {
		assert.Equal(t, "01:02:03", FormatTime(3723))
	})

	t.Run("Success - Large Value", func(t *testing.T) {
		assert.Equal(t, "10:00:00", FormatTime(36000))
	})

	t.Run("Success - Fractional Seconds Truncated", func(t *testing.T) {
		assert.Equal(t, "00:00:01", FormatTime(1.9))
	})

	t.Run("Success - Fractional Less Than One Second", func(t *testing.T) {
		assert.Equal(t, "00:00:00", FormatTime(0.5))
	})
}
