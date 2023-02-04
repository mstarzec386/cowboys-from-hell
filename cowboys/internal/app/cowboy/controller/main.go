package cowboyController

import (
	"github.com/gofiber/fiber/v2"

	// "cowboys/internal/pkg/clients/redis"
	"cowboys/internal/app/cowboy/state"
)

type CowboyController struct {
	GameState *state.GameState
}

func (g *CowboyController) Start(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func New() *CowboyController {
	gameState := state.New()

	return &CowboyController{GameState: gameState}
}
