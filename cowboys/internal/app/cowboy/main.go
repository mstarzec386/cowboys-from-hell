package cowboy

import (
	"fmt"

	"cowboys/internal/app/cowboy/server"
)


func Run(port int, gameMasterEdnpoint string) {
	fmt.Println("Cowboy Started")

	server.Run(port, gameMasterEdnpoint)
}