package limiter

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds the core configuration for the rate limiter.
// Framework-specific settings are handled in their respective middleware generators.
type Config struct {
	// Redis Configuration for starting limiter
	RedisClient *redis.Client
	RedisURL    string

	// Rate Limiter configuration
	MaxRequests int
	Window      time.Duration

	// value "token-bucket", "sliding-window" and "fixed-window"
	Algorithm string
}

type Limiter struct {
	store      Store
	config     Config
	ctx        context.Context
	cancelfunc context.CancelFunc
}

func New(config Config) (*Limiter, error) {
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize store
	store, err := initStore(ctx, config)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("store initialization failed: %w", err)
	}

	return &Limiter{
		store:      store,
		config:     config,
		ctx:        ctx,
		cancelfunc: cancel,
	}, nil
}

func (l *Limiter) Close() error {
	l.cancelfunc()

	if closer, ok := l.store.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

// Helper functions
func validateConfig(cfg *Config) error {
	if cfg.MaxRequests <= 0 {
		return errors.New("maxRequests must be positive")
	}
	if cfg.Window <= 0 {
		return errors.New("window duration must be positive")
	}
	if !slices.Contains([]string{"token-bucket", "sliding-window", "fixed-window"}, cfg.Algorithm) {
		return errors.New("invalid algorithm")
	}
	return nil
}

func initStore(ctx context.Context, config Config) (Store, error) {
	switch {
	case config.RedisClient != nil:
		return NewRedisStore(config.RedisClient), nil
	case config.RedisURL != "":
		rdb := redis.NewClient(&redis.Options{Addr: config.RedisURL})
		if err := rdb.Ping(ctx).Err(); err != nil {
			return nil, fmt.Errorf("redis connection failed: %w", err)
		}
		return NewRedisStore(rdb), nil
	default:
		return NewMemoryStore(), nil
	}
}
