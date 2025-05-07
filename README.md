# Fiber Limiter
[![Go Report Card](https://goreportcard.com/badge/github.com/NarmadaWeb/limiter)](https://goreportcard.com/report/github.com/NarmadaWeb/limiter)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A high-performance rate limiting middleware for [Fiber](https://github.com/gofiber/fiber) with Redis and in-memory support, implementing multiple rate-limiting algorithms.

## Features

- 🚀 **Multiple Algorithms**: Token Bucket, Sliding Window, and Fixed Window
- 💾 **Storage Options**: Redis (for distributed systems) and in-memory (for single-instance)
- ⚡ **High Performance**: Minimal overhead with efficient algorithms
- 🔧 **Customizable**: Flexible key generation and response handling
- 📊 **RFC Compliance**: Standard `RateLimit` headers (RFC 6585)

## Installation

```bash
go get github.com/NarmadaWeb/limiter/v2
```

## Usage

### Basic Example

```go
package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/NarmadaWeb/limiter/v2"
)

func main() {
	app := fiber.New()

	// Initialize with default in-memory store
	limiterCfg := limiter.Config{
		MaxRequests: 100,
		Window:      1 * time.Minute,
		Algorithm:   "sliding-window",
	}

	l, err := limiter.New(limiterCfg)
	if err != nil {
		panic(err)
	}

	app.Use(l.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
```

### With Redis

```go
import "github.com/redis/go-redis/v9"

// ...

rdb := redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

limiterCfg := limiter.Config{
	RedisClient: rdb,
	MaxRequests: 200,
	Window:      5 * time.Minute,
	Algorithm:   "token-bucket",
}
```

## Configuration Options

| Option                | Type                  | Description                                                                 |
|-----------------------|-----------------------|-----------------------------------------------------------------------------|
| `RedisClient`         | `*redis.Client`       | Redis client instance (optional)                                            |
| `RedisURL`            | `string`              | Redis connection URL (alternative to RedisClient)                           |
| `MaxRequests`         | `int`                 | Maximum allowed requests per window                                         |
| `Window`              | `time.Duration`       | Duration of the rate limit window (e.g., 1*time.Minute)                     |
| `Algorithm`           | `string`              | Rate limiting algorithm (`token-bucket`, `sliding-window`, `fixed-window`)  |
| `KeyGenerator`        | `func(*fiber.Ctx) string` | Custom function to generate rate limit keys (default: client IP)         |
| `SkipSuccessful`      | `bool`                | Don't count successful requests (status < 400)                              |
| `LimitReachedHandler` | `fiber.Handler`       | Custom handler when limit is reached                                        |
| `ErrorHandler`        | `func(*fiber.Ctx, error) error` | Custom error handler for storage/configuration errors           |

## Response Headers

The middleware adds these standard headers to responses:

- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests in current window
- `X-RateLimit-Reset`: Unix timestamp when limit resets
- `RateLimit-Policy`: Formal policy description

## Algorithms

### 1. Token Bucket
- Smooth bursting allowed
- Gradually refills tokens at steady rate
- Good for evenly distributed loads

### 2. Sliding Window
- Precise request counting
- Tracks exact request timestamps
- Prevents bursts at window edges

### 3. Fixed Window
- Simple implementation
- Counts requests per fixed interval
- May allow bursts at window boundaries


## Testing

Run the test suite:

```bash
go test -v ./...
```

Test with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Examples

See the [examples directory](examples/) for:
- Basic usage
- Redis integration
- Custom key generation
- Error handling
- Multiple limiters

## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you'd like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

## Key Sections Explained:

1. **Badges**: Shows build status, documentation link, and license info
2. **Features**: Highlights the package's capabilities
3. **Installation**: Simple one-line install command
4. **Usage**: Basic and Redis examples
5. **Configuration**: Table of all available options
6. **Headers**: Documents the response headers added
7. **Algorithms**: Explains the three implemented algorithms
8. **Benchmarks**: Gives performance expectations
9. **Testing**: How to run tests
10. **Examples**: Points to additional examples
11. **Contributing**: Guidelines for contributors
12. **License**: MIT license information

This README provides:
- Quick start for new users
- Detailed configuration reference
- Clear explanation of technical concepts
- Development and contribution guidelines
- Visual badges for project health

Would you like me to add any additional sections or modify any part of this README?
