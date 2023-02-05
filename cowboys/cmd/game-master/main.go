package main

import (
	"os"
	"strconv"

	"cowboys/internal/app/game-master"
)

func main() {
	gameMaster.Run(getPort())
}

func getPort() int {
	portStr := os.Getenv("GAME_MASTER_SERVICE_PORT")

	port, err := strconv.Atoi(portStr)

	if err != nil {
		port = 8000
	}

	return port
}