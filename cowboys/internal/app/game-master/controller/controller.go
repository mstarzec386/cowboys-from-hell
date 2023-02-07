package gameController

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	// "cowboys/internal/pkg/clients/redis"
	"cowboys/internal/app/game-master/state"
	"cowboys/internal/pkg/cowboys"
	"cowboys/internal/pkg/clients/redis"
)

type GameController struct {
	GameState *state.GameState
}

func (g *GameController) Status(c *fiber.Ctx) error {
	return c.JSON(g.GameState.Status)
}

func (g *GameController) Register(c *fiber.Ctx) error {
	registerData := new(cowboys.RegisterCowboy)

	if err := c.BodyParser(registerData); err != nil {
		return err
	}

	if err := validateRegisterBody(registerData); err != nil {
		return err
	}

	registerResponse := g.GameState.RegisterCowboy(registerData)

	if registerResponse == nil {
		return fiber.NewError(404, "No cowboys available")
	}

	fmt.Printf("Cowboy Registered\n%s\nEndpoint: %s\n", registerResponse.String(), registerData.ToUrl(""))

	return c.JSON(registerResponse)
}
func (g *GameController) Update(c *fiber.Ctx) error {
	updateData := new(cowboys.UpdateCowboy)
	id := c.Params("cowboyId")

	if err := c.BodyParser(updateData); err != nil {
		return err
	}

	if err := g.GameState.UpdateCowboy(id, updateData); err != nil {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func (g *GameController) GetAll(c *fiber.Ctx) error {
	return c.JSON(g.GameState.RegisteredPlayers)
}

func New(redisClient redis.RedisClientInterface) *GameController {
	gameState := state.New(redisClient)

	return &GameController{GameState: gameState}
}
