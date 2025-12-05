package limiter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEchoMiddleware(t *testing.T) {
	e := echo.New()

	cfg := limiter.Config{
		MaxRequests: 5,
		Window:      time.Minute,
		Algorithm:   "fixed-window",
	}
	l, err := limiter.New(cfg)
	assert.NoError(t, err)

	e.Use(l.EchoMiddleware(limiter.EchoConfig{}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}
