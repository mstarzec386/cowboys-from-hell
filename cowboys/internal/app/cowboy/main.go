
package cowboy

import (
	"github.com/gofiber/fiber/v2"

	"cowboys/internal/pkg/cowboys"
)

func Run() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

	app.Get("/status", func(c *fiber.Ctx) error {
		c.Accepts("application/json")

		r := cowboys.Cowboy{Name: "fuu", Health: 10, Damage: 1}

		return c.JSON(r)
	})

    app.Listen(":8000")
}