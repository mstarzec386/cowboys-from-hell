package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// TODO move to the cowboy repo as CowboyConfig or smth
type RegisterResponse struct {
	Name string `json:"name" xml:"name" form:"name"`
	Health int `json:"health" xml:"health" form:"health"`
	Damage int `json:"damage" xml:"damage" form:"damage"`
}

type RegisterRequestBody struct {
	Host string `json:"host" xml:"host" form:"host"`
	Port int `json:"port" xml:"port" form:"port"`
}

func Run() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Get("/status", func(c *fiber.Ctx) error {
        return c.SendString("ixi")
    })

	app.Get("/start", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		body := new(RegisterRequestBody)

		if err := c.BodyParser(body); err != nil {
            return err
        }

		fmt.Printf("Req %v\n", body)

		// TODO get from DB, randomly assign
		r := RegisterResponse{Name: "fuu", Health: 10, Damage: 1}

		return c.JSON(r)
	})

    app.Listen(":8000")

	// TODO
	//Register()
}