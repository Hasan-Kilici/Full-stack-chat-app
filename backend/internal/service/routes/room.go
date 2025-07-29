package routes

import (
	"github.com/hasan-kilici/chat/pkg/middleware"
	"github.com/hasan-kilici/chat/internal/service/handlers"
	"github.com/gofiber/fiber/v3"
)

func Room(app fiber.Router) {
	app.Post("/room",middleware.JWTAuth, handlers.CreateRoomHandler)
	app.Get("/room/:room_id/join", middleware.JWTAuth, handlers.JoinRoomHandler)
}