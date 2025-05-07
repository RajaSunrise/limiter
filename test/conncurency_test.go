package limiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentAccess(t *testing.T) {
	store := limiter.NewMemoryStore()
	ctx := context.Background()
	key := "concurrent-test"
	maxRequests := 100
	var wg sync.WaitGroup
	var allowedCount int

	for i := 0; i < maxRequests*2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			allowed, _, _, err := store.Take(ctx, key, maxRequests, time.Minute, "fixed-window")
			assert.NoError(t, err)
			if allowed {
				allowedCount++
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, maxRequests, allowedCount)
}
