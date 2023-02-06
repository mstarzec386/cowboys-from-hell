package gameMaster

import (
	"fmt"

	"cowboys/internal/app/game-master/server"
)

func Run(port int, redisHost string) {
	fmt.Println("Game Server started")

	server.Run(port, redisHost)
}
