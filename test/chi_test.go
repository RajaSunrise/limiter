package limiter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestChiMiddleware(t *testing.T) {
	r := chi.NewRouter()

	cfg := limiter.Config{
		MaxRequests: 5,
		Window:      time.Minute,
		Algorithm:   "fixed-window",
	}
	l, err := limiter.New(cfg)
	assert.NoError(t, err)

	// StdLibMiddleware works for Chi
	r.Use(l.StdLibMiddleware(limiter.StdLibConfig{}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
