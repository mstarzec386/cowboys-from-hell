package cowboy

import (
	"fmt"

	"cowboys/internal/app/cowboy/server"
)


func Run(port int) {
	fmt.Println("Cowboy Started")

	server.Run(port)
}