package routes

import (
	"github.com/hasan-kilici/chat/pkg/middleware"
	"github.com/hasan-kilici/chat/internal/service/handlers"
	"github.com/gofiber/fiber/v3"
)

func Api(app fiber.Router) {
	app.Get("/search/user", middleware.JWTAuth, handlers.SearchUser)
	app.Get("/get/messages/:room_id", middleware.JWTAuth, handlers.GetMessages)
	app.Get("/get/participants/", middleware.JWTAuth, handlers.GetUserRooms)
}