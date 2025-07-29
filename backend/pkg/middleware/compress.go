package middleware

import(
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3"
)

// Compress is a middleware that compresses HTTP responses.
func Compress(c fiber.Ctx) error {
	compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
	return c.Next()
}