package queue_test

import (
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/queue"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAsynqClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mr, err := miniredis.Run()
		require.NoError(t, err)
		defer mr.Close()

		client, err := queue.NewAsynqClient(mr.Addr())

		assert.NoError(t, err)
		assert.NotNil(t, client)
	})

	t.Run("Error - URL", func(t *testing.T) {
		client, err := queue.NewAsynqClient("")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "Redis address is empty")
		assert.Nil(t, client)
	})

	t.Run("Error - Ping", func(t *testing.T) {
		client, err := queue.NewAsynqClient("redis: //localhost:1")

		assert.Error(t, err)
		assert.Nil(t, client)
	})
}
