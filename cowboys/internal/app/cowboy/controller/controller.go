package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"cowboys/internal/app/cowboy/game"
)

type CowboyController struct {
	Game *game.Game
}

func (g *CowboyController) Start(c *fiber.Ctx) error {
	if !g.Game.IsRegistered() {
		c.Status(400)
		return c.SendString("Cowboy not registered")
	}

	if err := g.Game.Start(); err != nil {
		c.Status(400)
		return c.SendString(err.Error())
	}

	return c.SendStatus(200)
}

func (g *CowboyController) Stop(c *fiber.Ctx) error {
	if err := g.Game.Stop(); err != nil {
		c.Status(400)
		return c.SendString(err.Error())
	}

	return c.SendStatus(200)
}

func (g *CowboyController) Stats(c *fiber.Ctx) error {
	return c.JSON(g.Game.GetCowboy())
}

func (g *CowboyController) Hit(c *fiber.Ctx) error {
	damageString := c.Params("damage")

	damage, err := strconv.Atoi(damageString)

	if err != nil {
		c.Status(400)
		return c.SendString("Damage not a number")
	}

	if health := g.Game.GetHealth(); health < 1 {
		c.Status(400)
		return c.SendString("Already dead")
	}

	g.Game.HitCowboy(damage)

	return c.SendStatus(200)
}

func New(game *game.Game) *CowboyController {
	return &CowboyController{Game: game}
}
