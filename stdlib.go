package limiter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StdLibConfig struct {
	KeyGenerator        func(r *http.Request) string
	LimitReachedHandler http.HandlerFunc
	ErrorHandler        func(w http.ResponseWriter, r *http.Request, err error)
	Skipsuccessfull     bool
}

// StdLibMiddleware creates a standard net/http middleware.
// This works for Chi, Go Standard Library, and any framework compatible with http.Handler.
func (l *Limiter) StdLibMiddleware(cfg StdLibConfig) func(http.Handler) http.Handler {
	// Set defaults
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = func(r *http.Request) string {
			// Basic IP extraction
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
				// Remove port if present
				if idx := strings.LastIndex(ip, ":"); idx != -1 {
					ip = ip[:idx]
				}
			}
			return ip
		}
	}
	if cfg.LimitReachedHandler == nil {
		cfg.LimitReachedHandler = func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "rate limit exceeded",
				"message": "Too many requests, please try again later",
			})
		}
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "rate limit error",
				"message": err.Error(),
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := cfg.KeyGenerator(r)

			allowed, remaining, reset, err := l.store.Take(l.ctx, key, l.config.MaxRequests, l.config.Window, l.config.Algorithm)
			if err != nil {
				cfg.ErrorHandler(w, r, err)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.config.MaxRequests))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(reset.Unix(), 10))
			w.Header().Set("RateLimit-Policy", fmt.Sprintf("%d;w=%d", l.config.MaxRequests, int(time.Minute.Seconds())))

			if !allowed {
				cfg.LimitReachedHandler(w, r)
				return
			}

			// To handle Skipsuccessfull, we need to capture the status code.
			// Wrap ResponseWriter
			ww := &responseWriter{ResponseWriter: w, code: http.StatusOK}
			next.ServeHTTP(ww, r)

			if cfg.Skipsuccessfull && ww.code < http.StatusBadRequest {
				_ = l.store.Rollback(l.ctx, key)
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
	rw.ResponseWriter.WriteHeader(code)
}
