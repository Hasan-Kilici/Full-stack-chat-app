package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

// RateLimit is a middleware that limits the number of requests.
var RateLimit = limiter.New(limiter.Config{
	Max:        20,
	Expiration: 1 * time.Minute,
	KeyGenerator: func(c fiber.Ctx) string {
		return c.IP()
	},
	LimitReached: func(c fiber.Ctx) error {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Rate limit exceeded. Please try again later.",
		})
	},
})