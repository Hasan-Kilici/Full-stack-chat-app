package service

import (
	"github.com/hasan-kilici/chat/internal/service/repository"
	"github.com/hasan-kilici/chat/internal/service/routes"
	"github.com/hasan-kilici/chat/pkg/middleware"
	//"github.com/hasan-kilici/chat/pkg/utils"
	
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func Start() error {
	//Config := utils.LoadConfig("./configs/service.ini")
	repository.Connect()
	app := fiber.New()

	app.Use(middleware.Cors)
	app.Use(middleware.Logger)
	app.Use(middleware.Compress)
	app.Use(middleware.Security)
	// app.Use(middleware.RateLimit)
	app.Use(recover.New())

	//service := app.Group("/")
	auth := app.Group("/auth")
	room := app.Group("/room")
	api := app.Group("/api")

	routes.Api(api)
	routes.Auth(auth)
	routes.Room(room)
	
	app.Use(middleware.NotFound)

	err := app.Listen("127.0.0.1:4000")
	if err != nil {
		panic(err)
	}

	return err
}