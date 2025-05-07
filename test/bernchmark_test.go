package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
)

func BenchmarkLimiter(b *testing.B) {
	store := limiter.NewMemoryStore()
	ctx := context.Background()
	key := "benchmark-test"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _ = store.Take(ctx, key, 1000, time.Minute, "fixed-window")
	}
}
