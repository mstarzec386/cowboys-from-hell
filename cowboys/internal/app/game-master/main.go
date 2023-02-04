package gameMaster

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"cowboys/internal/pkg/cowboys"
)


func Run() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

	app.Post("/register", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		body := new(cowboys.RegisterRequestBody)

		if err := c.BodyParser(body); err != nil {
            return err
        }

		if err := validateRegister(body); err != nil {
			return err
		}

		fmt.Printf("Req %v\n", body)

		// TODO get from DB, randomly assign
		r := cowboys.Cowboy{Name: "fuu", Health: 10, Damage: 1}

		return c.JSON(r)
	})

    app.Listen(":8000")
}