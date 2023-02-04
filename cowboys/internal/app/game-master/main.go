package gameMaster

import (
	"fmt"

	"cowboys/internal/app/game-master/server"
)


func Run(port int) {
	fmt.Println("Game Server started")

	server.Run(port)
}