package limiter

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type GinConfig struct {
	KeyGenerator        func(c *gin.Context) string
	LimitReachedHandler func(c *gin.Context)
	ErrorHandler        func(c *gin.Context, err error)
	Skipsuccessfull     bool
}

func (l *Limiter) GinMiddleware(cfg GinConfig) gin.HandlerFunc {
	// Set defaults
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = func(c *gin.Context) string {
			return c.ClientIP()
		}
	}
	if cfg.LimitReachedHandler == nil {
		cfg.LimitReachedHandler = func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c *gin.Context, err error) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "rate limit error",
				"message": err.Error(),
			})
			c.Abort()
		}
	}

	return func(c *gin.Context) {
		key := cfg.KeyGenerator(c)

		allowed, remaining, reset, err := l.store.Take(l.ctx, key, l.config.MaxRequests, l.config.Window, l.config.Algorithm)
		if err != nil {
			cfg.ErrorHandler(c, err)
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(l.config.MaxRequests))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
		c.Header("RateLimit-Policy", fmt.Sprintf("%d;w=%d", l.config.MaxRequests, int(time.Minute.Seconds())))

		if !allowed {
			cfg.LimitReachedHandler(c)
			return
		}

		c.Next()

		if cfg.Skipsuccessfull && c.Writer.Status() < http.StatusBadRequest {
			_ = l.store.Rollback(l.ctx, key)
		}
	}
}
