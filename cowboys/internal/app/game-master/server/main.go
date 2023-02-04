package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"cowboys/internal/app/game-master/controller"
)

func Run(port int) {
	app := fiber.New()
	game := gameController.New()

	app.Get("/status", game.Status)
	app.Post("/register", game.Register)

	app.Listen(fmt.Sprintf(":%d", port))
}
