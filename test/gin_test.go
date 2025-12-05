package limiter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	cfg := limiter.Config{
		MaxRequests: 5,
		Window:      time.Minute,
		Algorithm:   "fixed-window",
	}
	l, err := limiter.New(cfg)
	assert.NoError(t, err)

	router.Use(l.GinMiddleware(limiter.GinConfig{}))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
