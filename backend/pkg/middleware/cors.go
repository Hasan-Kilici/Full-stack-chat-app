package middleware

import "github.com/gofiber/fiber/v3/middleware/cors"

// Cors is a middleware that handles CORS.
var Cors = cors.New(cors.Config{
    AllowOrigins: []string{"*"},
    AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
})