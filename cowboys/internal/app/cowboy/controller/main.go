package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"cowboys/internal/app/cowboy/state"
)

type CowboyController struct {
	GameState *state.GameState
}

func (g *CowboyController) Start(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func (g *CowboyController) Stats(c *fiber.Ctx) error {
	return c.JSON(g.GameState.GetCowboy())
}

func (g *CowboyController) Hit(c *fiber.Ctx) error {
	damageString := c.Params("damage")

	damage, err := strconv.Atoi(damageString)

	if err != nil {
		c.Status(400)
		return c.SendString("Damage not a number")
	}

	if health := g.GameState.GetHealth(); health < 1 {
		c.Status(400)
		return c.SendString("Already dead")
	}

	g.GameState.HitCowboy(damage)

	return c.SendStatus(200)
}

func New(state *state.GameState) *CowboyController {
	return &CowboyController{GameState: state}
}
