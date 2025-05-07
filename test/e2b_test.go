package limiter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/NarmadaWeb/limiter/v2"
	"github.com/stretchr/testify/assert"
)

func TestEndToEnd(t *testing.T) {
	app := fiber.New()

	limiterCfg := limiter.Config{
		MaxRequests: 2,
		Window:      1 * time.Second,
		Algorithm:   "fixed-window",
	}
	l, err := limiter.New(limiterCfg)
	assert.NoError(t, err)

	app.Use(l.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Test helper function
	makeRequest := func() int {
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		return resp.StatusCode
	}

	// First two requests should succeed
	assert.Equal(t, http.StatusOK, makeRequest())
	assert.Equal(t, http.StatusOK, makeRequest())

	// Third request should fail
	assert.Equal(t, http.StatusTooManyRequests, makeRequest())

	// Wait for window to reset
	time.Sleep(1100 * time.Millisecond) // Slightly more than the window

	// After reset, next request should succeed
	assert.Equal(t, http.StatusOK, makeRequest())
}
