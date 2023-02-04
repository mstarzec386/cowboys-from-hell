package gameMaster

import (
	"cowboys/internal/pkg/cowboys"
	"errors"
)

func validateRegister(body *cowboys.RegisterRequestBody) error {
	if (body.Host == "") {
		return errors.New("body: Empty host")
	}

	if (body.Port == 0) {
		return errors.New("body: Empty Port")
	}

	return nil
}