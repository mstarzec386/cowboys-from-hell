package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"cowboys/internal/app/game-master/controller"
)

func Run(port int) {
	app := fiber.New()

	game := gameController.New()

	app.Use(logger.New(logger.Config{
		Format: "Request from [${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/status", game.Status)
	// TODO not necessary but should be ;)
	// app.Get("/cowboys/:cowboyId", game.Get)
	app.Get("/cowboys", game.GetAll)
	app.Post("/cowboys", game.Register)
	app.Put("/cowboys/:cowboyId", game.Update)

	app.Listen(fmt.Sprintf(":%d", port))
}
