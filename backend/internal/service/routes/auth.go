package routes
import (
	"github.com/hasan-kilici/chat/pkg/middleware"
	"github.com/hasan-kilici/chat/internal/service/handlers"
	"github.com/gofiber/fiber/v3"
)

func Auth(app fiber.Router) {
	app.Post("/register", handlers.RegisterHandler)
	app.Post("/login", handlers.LoginHandler)
	app.Get("/profile", middleware.JWTAuth, handlers.AuthRequired)
	app.Get("/@me", middleware.JWTAuth, handlers.GetUserProfile)
}