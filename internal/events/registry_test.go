package events_test

import (
	"context"
	"testing"

	"github.com/CABGenOrg/cabgen_backend/internal/events"
	"github.com/stretchr/testify/assert"
)

func TestNewRegistry(t *testing.T) {
	result := events.NewRegistry()

	assert.NotNil(t, result)
}

func TestRegister(t *testing.T) {
	name := "user.created"
	function := func(ctx context.Context, payload []byte) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {
		reg := events.NewRegistry()
		err := reg.Register(name, function)

		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		reg := events.NewRegistry()
		err := reg.Register(name, function)
		assert.NoError(t, err)

		err = reg.Register(name, function)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "handler already registered for event")
	})
}

func TestGet(t *testing.T) {
	name := "user.created"
	function := func(ctx context.Context, payload []byte) error {
		return nil
	}

	t.Run("Success", func(t *testing.T) {
		reg := events.NewRegistry()
		err := reg.Register(name, function)
		assert.NoError(t, err)

		result, ok := reg.Get(name)
		assert.NotNil(t, result)
		assert.True(t, ok)
	})

	t.Run("Error", func(t *testing.T) {
		reg := events.NewRegistry()
		result, ok := reg.Get(name)

		assert.Nil(t, result)
		assert.False(t, ok)
	})
}
