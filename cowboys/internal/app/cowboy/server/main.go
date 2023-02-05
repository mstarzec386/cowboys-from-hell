package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"cowboys/internal/app/cowboy/controller"
	"cowboys/internal/app/cowboy/game"
	"cowboys/internal/app/cowboy/state"
)

func Run(port int, gameMasterEndpoint string) {
	app := fiber.New()

	gameState := state.New()
	cowboyGame := game.New(gameState, gameMasterEndpoint)
	cowboyController := controller.New(cowboyGame)

	app.Use(logger.New(logger.Config{
		Format: "Request from [${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/stats", cowboyController.Stats)
	app.Get("/start", cowboyController.Start)
	app.Get("/stop", cowboyController.Stop)
	app.Get("/hit/:damage", cowboyController.Hit)

	go cowboyGame.Register()

	app.Listen(fmt.Sprintf(":%d", port))
}
