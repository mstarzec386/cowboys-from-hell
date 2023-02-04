package gameController

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	// "cowboys/internal/pkg/clients/redis"
	"cowboys/internal/pkg/cowboys"
)

type GameCowboy struct {
	Endpoint *cowboys.RegisterRequestBody `json:"endpoint" xml:"endpoint" form:"endpoint"`
	Cowboy   *cowboys.Cowboy              `json:"cowboy" xml:"cowboy" form:"cowboy"`
}

func (c GameCowboy) String() string {
	return fmt.Sprintf("Name: %s, Health: %d, Damage: %d, Host: %s, Port %d",
		c.Cowboy.Name, c.Cowboy.Health, c.Cowboy.Damage, c.Endpoint.Host, c.Endpoint.Port)
}

type GameState struct {
	RegisteredPlayers []GameCowboy `json:"registeredPlayers" xml:"registeredPlayers" form:"registeredPlayers"`
	Status            string       `json:"status" xml:"status" form:"status"`
}

type GameController struct {
	// redisClient    redisClient.RedisClient
	gameState      GameState
	initialPlayers []*cowboys.Cowboy
	playersNumbers int
}

func (g *GameController) Status(c *fiber.Ctx) error {
	return c.JSON(g.gameState)
}

func (g *GameController) Register(c *fiber.Ctx) error {
	c.Accepts("application/json")
	body := new(cowboys.RegisterRequestBody)

	if err := c.BodyParser(body); err != nil {
		return err
	}

	if err := validateRegisterBody(body); err != nil {
		return err
	}

	newCowboy := g.registerCowboy()
	if newCowboy == nil {
		return fiber.NewError(404, "No cowboys available")
	}

	cowboy := GameCowboy{Cowboy: newCowboy, Endpoint: body}
	g.gameState.RegisteredPlayers = append(g.gameState.RegisteredPlayers, cowboy)

	fmt.Printf("Register Cowboy: %s\n", cowboy.String())

	return c.JSON(newCowboy)
}

func (g *GameController) registerCowboy() *cowboys.Cowboy {
	lastOneIndex := len(g.initialPlayers)

	if lastOneIndex > 0 {
		cowboy := g.initialPlayers[lastOneIndex-1]
		g.initialPlayers = g.initialPlayers[:lastOneIndex-1]

		return cowboy
	}

	return nil
}

func New() *GameController {
	// TODO get players from redis
	var initialPlayers []*cowboys.Cowboy
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Eliot", Health: 10, Damage: 1})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Fliz", Health: 10, Damage: 2})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Pawl", Health: 5, Damage: 1})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Dvil", Health: 15, Damage: 3})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Gatt", Health: 6, Damage: 1})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Luci", Health: 12, Damage: 2})
	playerNumbers := len(initialPlayers)

	gameState := GameState{Status: "Register"}

	return &GameController{initialPlayers: initialPlayers, playersNumbers: playerNumbers, gameState: gameState}
}
