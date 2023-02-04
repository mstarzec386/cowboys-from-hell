package gameController

import (
	"github.com/gofiber/fiber/v2"

	"cowboys/internal/pkg/cowboys"
)

func validateRegisterBody(body *cowboys.RegisterRequestBody) error {
	if body.Host == "" {
		return fiber.NewError(400, "Empty host")
	}

	if body.Port == 0 {
		return fiber.NewError(400, "Empty port")
	}

	return nil
}
