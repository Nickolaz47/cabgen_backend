package redis_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/CABGenOrg/cabgen_backend/internal/config"
	redisclient "github.com/CABGenOrg/cabgen_backend/internal/redis"
)

func TestNewRedisClient_Success(t *testing.T) {
	mr, err := miniredis.Run()
	require.NoError(t, err)
	defer mr.Close()

	oldRedisURL := config.RedisURL
	config.RedisURL = "redis://" + mr.Addr()
	defer func() {
		config.RedisURL = oldRedisURL
	}()

	client, err := redisclient.NewRedisClient()

	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestNewRedisClient_InvalidURL(t *testing.T) {
	oldRedisURL := config.RedisURL
	config.RedisURL = "not-a-valid-url"
	defer func() {
		config.RedisURL = oldRedisURL
	}()

	_, err := redisclient.NewRedisClient()

	assert.ErrorContains(t, err, "Failed to parse Redis URL")
}

func TestNewRedisClient_ConnectionRefused(t *testing.T) {
	oldRedisURL := config.RedisURL
	config.RedisURL = "redis://localhost:1"
	defer func() {
		config.RedisURL = oldRedisURL
	}()
	_, err := redisclient.NewRedisClient()

	assert.ErrorContains(t, err, "Failed to connect to Redis")
}
