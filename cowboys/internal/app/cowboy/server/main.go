package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"cowboys/internal/app/cowboy/controller"
)

func Run(port int) {
	app := fiber.New()

	cowboy := cowboyController.New()

	app.Use(logger.New(logger.Config{
		Format: "Request from [${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/start", cowboy.Start)

	app.Listen(fmt.Sprintf(":%d", port))
}
