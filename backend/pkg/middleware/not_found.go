package middleware

import "github.com/gofiber/fiber/v3"

// NotFound handles 404 errors.
func NotFound(c fiber.Ctx) error {
	c.Status(404).JSON(fiber.Map{
		"code":    404,
		"message": "404: Not Found",
	})

	return nil
}