package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/NarmadaWeb/limiter/v2"
)

func main() {
	app := fiber.New()

	cfg := limiter.Config{
		MaxRequests: 3,
		Window:      30 * time.Second,
		Algorithm: "sliding-window",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.Printf("Rate limiter error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "rate_limit_error",
				"message": "Unable to process rate limit",
			})
		},
		LimitReachedHandler: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "rate_limit_exceeded",
				"message": "Slow down! You've made too many requests",
				"retry":   c.Get("X-RateLimit-Reset"),
			})
		},
	}

	l, err := limiter.New(cfg)
	if err != nil {
		panic(err)
	}

	app.Use(l.Middleware())

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("API endpoint with custom error handling")
	})

	app.Listen(":3000")
}
