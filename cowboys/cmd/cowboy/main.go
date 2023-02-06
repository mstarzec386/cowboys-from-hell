package main

import (
	"cowboys/internal/app/cowboy"
)

func main() {
	// TODO config would be a nice thing to add or cmd arguments/envs
	cowboy.Run(8000, "http://game-master:8000")
}