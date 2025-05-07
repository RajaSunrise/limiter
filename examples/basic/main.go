package main

import (
	"log"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize with default in-memory store
	cfg := limiter.Config{
		MaxRequests: 5,
		Window:      1 * time.Minute,
		Algorithm:   "fixed-window",
	}

	l, err := limiter.New(cfg)
	if err != nil {
		panic(err)
	}

	app.Use(l.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":   "Hello World",
			"remaining": c.Get("X-RateLimit-Remaining"),
		})
	})

	log.Fatal(app.Listen(":3000"))
}
