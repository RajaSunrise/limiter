package main

import (
	"log"
	"time"

	"github.com/NarmadaWeb/limiter/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	cfg := limiter.Config{
		MaxRequests: 10,
		Window:      1 * time.Hour,
		Algorithm:   "token-bucket",
		KeyGenerator: func(c *fiber.Ctx) string {
			if apiKey := c.Get("X-API-Key"); apiKey != "" {
				return "api:" + apiKey
			}
			return "ip:" + c.IP() + ":" + c.Path()
		},
	}

	l, err := limiter.New(cfg)
	if err != nil {
		panic(err)
	}

	app.Use(l.Middleware())

	app.Get("/profile", func(c *fiber.Ctx) error {
		return c.SendString("Profile page - custom key rate limiting")
	})

	log.Fatal(app.Listen(":3000"))
}
