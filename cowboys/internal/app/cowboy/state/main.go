package state

import (
	"sync"

	"cowboys/internal/pkg/cowboys"
)

type GameState struct {
	mu     sync.Mutex
	id     string
	cowboy *cowboys.Cowboy
}

func (g *GameState) GetCowboy() cowboys.Cowboy {
	g.mu.Lock()
	defer g.mu.Unlock()

	return *g.cowboy
}

func (g *GameState) GetHealth() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.cowboy.Health
}

func (g *GameState) SetId(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.id = id
}

func (g *GameState) GetId() string {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.id
}

func (g *GameState) SetCowboy(cowboy cowboys.Cowboy) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.cowboy = &cowboy
}

func (g *GameState) HitCowboy(damage int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.cowboy.Health < damage {
		g.cowboy.Health = 0

	} else {
		g.cowboy.Health -= damage
	}
}

func New() *GameState {
	return &GameState{cowboy: &cowboys.Cowboy{}}
}
