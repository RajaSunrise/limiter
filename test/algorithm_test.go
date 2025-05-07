package limiter_test

import (
	"context"
	"testing"
	"time"
	"github.com/NarmadaWeb/limiter/v2"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucketAlgorithm(t *testing.T) {
	store := limiter.NewMemoryStore()
	ctx := context.Background()
	key := "token-bucket-test"

	// Test initial request
	allowed, remaining, _, err := store.Take(ctx, key, 10, time.Minute, "token-bucket")
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, 9, remaining)

	// Test burst
	for i := 0; i < 9; i++ {
		_, _, _, err = store.Take(ctx, key, 10, time.Minute, "token-bucket")
		assert.NoError(t, err)
	}

	// Test limit exceeded
	allowed, _, _, err = store.Take(ctx, key, 10, time.Minute, "token-bucket")
	assert.NoError(t, err)
	assert.False(t, allowed)
}

func TestSlidingWindowAlgorithm(t *testing.T) {
	store := limiter.NewMemoryStore()
	ctx := context.Background()
	key := "sliding-window-test"

	// Test initial requests
	for i := 0; i < 5; i++ {
		allowed, remaining, _, err := store.Take(ctx, key, 5, time.Minute, "sliding-window")
		assert.NoError(t, err)
		assert.True(t, allowed)
		assert.Equal(t, 4-i, remaining)
	}

	// Test limit exceeded
	allowed, _, _, err := store.Take(ctx, key, 5, time.Minute, "sliding-window")
	assert.NoError(t, err)
	assert.False(t, allowed)
}
