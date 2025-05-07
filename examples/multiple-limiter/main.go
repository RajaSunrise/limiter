package main

import (
	"log"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Global rate limiter (10 requests/minute)
	globalLimiter, _ := limiter.New(limiter.Config{
		MaxRequests: 10,
		Window:      1 * time.Minute,
		Algorithm: 	 "sliding-window",
	})

	// Strict API limiter (2 requests/second)
	apiLimiter, _ := limiter.New(limiter.Config{
		MaxRequests: 2,
		Window:      1 * time.Second,
		Algorithm:   "sliding-window",
	})

	// Apply global limiter to all routes
	app.Use(globalLimiter.Middleware())

	// Public route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Public endpoint (global rate limit only)")
	})

	// API group with additional rate limiting
	api := app.Group("/api", apiLimiter.Middleware())
	api.Get("/data", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data":      "Sensitive API data",
			"remaining": c.Get("X-RateLimit-Remaining"),
		})
	})

	log.Fatal(app.Listen(":3000"))
}
