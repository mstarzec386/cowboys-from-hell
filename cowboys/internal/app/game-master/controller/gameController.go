package gameController

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	// "cowboys/internal/pkg/clients/redis"
	"cowboys/internal/pkg/cowboys"
)

const (
	Register   = "Register"
	Ready      = "Ready"
	InProgress = "In Progress"
	Done       = "Done"
)

type GameCowboy struct {
	Id       string                       `json:"id" xml:"id" form:"id"`
	Endpoint *cowboys.RegisterRequestBody `json:"endpoint" xml:"endpoint" form:"endpoint"`
	Cowboy   *cowboys.Cowboy              `json:"cowboy" xml:"cowboy" form:"cowboy"`
}

func (c GameCowboy) String() string {
	return fmt.Sprintf("Id: %s, Name: %s, Health: %d, Damage: %d, Host: %s, Port %d",
		c.Id, c.Cowboy.Name, c.Cowboy.Health, c.Cowboy.Damage, c.Endpoint.Host, c.Endpoint.Port)
}

// TODO move Game state to different module
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
	return c.JSON(g.gameState.Status)
}

func (g *GameController) Register(c *fiber.Ctx) error {
	registerData := new(cowboys.RegisterRequestBody)

	if err := c.BodyParser(registerData); err != nil {
		return err
	}

	if err := validateRegisterBody(registerData); err != nil {
		return err
	}

	registerResponse := g.registerCowboy(registerData)
	if registerResponse == nil {
		return fiber.NewError(404, "No cowboys available")
	}

	return c.JSON(registerResponse)
}
func (g *GameController) Update(c *fiber.Ctx) error {
	updateData := new(cowboys.UpdateRequestBody)
	id := c.Params("cowboyId")

	if err := c.BodyParser(updateData); err != nil {
		return err
	}

	for _, cowboy := range g.gameState.RegisteredPlayers {
		if cowboy.Id == id {
			cowboy.Cowboy.Health = updateData.Health
			return c.SendStatus(200)
		}
	}

	return c.SendStatus(404)
}

func (g *GameController) GetAll(c *fiber.Ctx) error {
	return c.JSON(g.gameState.RegisteredPlayers)
}

func (g *GameController) registerCowboy(registerData *cowboys.RegisterRequestBody) *cowboys.RegisterResponseBody {
	lastOneIndex := len(g.initialPlayers)

	if lastOneIndex > 0 {
		newCowboy := g.initialPlayers[lastOneIndex-1]
		g.initialPlayers = g.initialPlayers[:lastOneIndex-1]

		gameCowboy := GameCowboy{Cowboy: newCowboy, Endpoint: registerData, Id: generateId(newCowboy, registerData)}
		g.gameState.RegisteredPlayers = append(g.gameState.RegisteredPlayers, gameCowboy)
		// TODO  aadd g.mapCowboys id -> &GameCowboy

		fmt.Printf("Register Cowboy: %s\n", gameCowboy.String())

		if len(g.gameState.RegisteredPlayers) == g.playersNumbers {
			g.setInprogressStatus()
			go g.notifyCowboys()
		}

		return &cowboys.RegisterResponseBody{Id: gameCowboy.Id, Cowboy: newCowboy}
	}

	return nil
}

func (g *GameController) setInprogressStatus() {
	g.setStatus(Ready)
}

func (g *GameController) notifyCowboys() {
	for _, cowboy := range g.gameState.RegisteredPlayers {
		// TODO error handling wait for all responses etc
		go notifyCowboy(cowboy)
	}
}

func (g *GameController) setStatus(status string) {
	g.gameState.Status = status
}

func notifyCowboy(cowboy GameCowboy) {
	cowboyUrl := cowboy.Endpoint.ToUrl("start")
	resp, err := http.Get(cowboyUrl)
	if err != nil || resp.StatusCode != 200 {
		// TODO :D
		panic(err)
	}

	fmt.Printf("Cowboy notified: %s", cowboy.String())
}

func generateId(cowboy *cowboys.Cowboy, registerData *cowboys.RegisterRequestBody) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s-%s-%d", cowboy.Name, registerData.Host, registerData.Port)))

	return hex.EncodeToString(hash[:])
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

	gameState := GameState{Status: Register}

	return &GameController{initialPlayers: initialPlayers, playersNumbers: playerNumbers, gameState: gameState}
}
