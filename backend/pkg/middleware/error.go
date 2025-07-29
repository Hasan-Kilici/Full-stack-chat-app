package middleware

import "github.com/gofiber/fiber/v3"

func Error(c fiber.Ctx) error {
	c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    fiber.StatusInternalServerError,
		"message": "500: Internal server error",
	})

	return nil
}