package middleware

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/hasan-kilici/chat/internal/service/repository"
	"github.com/hasan-kilici/chat/pkg/auth"
)

func JWTAuth(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header format must be Bearer {token}",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Token is missing",
		})
	}

	claims, err := auth.ParseJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	session, err := repository.GetSessionByID(claims.SessionID)
	if err != nil || session == nil || session.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Session expired or invalid",
	    })
	}

	c.Locals("user_id", claims.UserID.String())
	c.Locals("session_id", claims.SessionID.String())

	return c.Next()
}
