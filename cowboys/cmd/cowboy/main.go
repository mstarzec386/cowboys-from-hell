package main

import (
	// "os"
	// "strconv"

	"cowboys/internal/app/cowboy"
)

func main() {
	cowboy.Run(8000, "http://game-master:8000")
}

/* func getPort() int {
	portStr := os.Getenv("COWBOYS_SERVICE_PORT")

	port, err := strconv.Atoi(portStr)

	if err != nil {
		port = 8000
	}

	return port
} */