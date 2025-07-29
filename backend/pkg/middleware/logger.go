package middleware

import "github.com/gofiber/fiber/v3/middleware/logger"

// Logger is a middleware that logs HTTP requests.
var Logger = logger.New(logger.Config{
    Format:     "${time} | ${pid} | ${latency} | ${status} - ${method} ${path} | ${ip} \n",
    TimeFormat: "02.01.2006 15:04:05",
})