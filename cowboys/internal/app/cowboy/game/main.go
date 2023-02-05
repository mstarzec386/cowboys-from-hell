package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"cowboys/internal/app/cowboy/state"
	"cowboys/internal/pkg/cowboys"
)

type Game struct {
	GameMasterEdnpoint string
	GameState          *state.GameState
	Id                 string
	mu                 sync.Mutex
	ticker             *time.Ticker
	done               chan bool
}

func (g *Game) GetGameState() *state.GameState {
	return g.GameState
}

func (g *Game) GetCowboy() cowboys.Cowboy {
	return g.GameState.GetCowboy()
}

func (g *Game) GetHealth() int {
	return g.GameState.GetHealth()
}

func (g *Game) HitCowboy(damage int) {
	g.GameState.HitCowboy(damage)
}

func (g *Game) Start() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.ticker != nil {
		return errors.New("game already started")
	}

	go g.runLoop()

	return nil
}

func (g *Game) runLoop() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	ticker := time.NewTicker(1 * time.Second)
	g.ticker = ticker

	// Use a channel to receive the Ticker's events
	done := make(chan bool)
	g.done = done

	go func() {
		for {
			select {
			case <-ticker.C:
				// Execute the "someFunction" function
				g.letsPlay()
			case <-done:
				// Stop the Ticker and exit the goroutine
				ticker.Stop()
				return
			}
		}
	}()

	// Stop the Ticker after some time

	return nil
}

func (g *Game) Stop() error {
	if g.ticker == nil {

		return errors.New("game already stopped")
	}
	g.done <- true
	g.ticker = nil

	return nil
}

func (g *Game) letsPlay() {
	fmt.Println("next round")

}

func (g *Game) IsRegistered() bool {
	if g.GameState.GetId() != "" {
		return true
	} else {
		return false
	}
}

func (g *Game) Register() error {
	fmt.Println("register")

	registerEndpoint := fmt.Sprintf("%s/cowboys", g.GameMasterEdnpoint)
	myEndpoint := getMyEndpoint()
	myEndpointJson, err := json.Marshal(myEndpoint)
	if err != nil {
		return err
	}

	fmt.Printf("register: %s with %s", registerEndpoint, myEndpointJson)

	res, err := http.Post(registerEndpoint, "application/json", bytes.NewBuffer([]byte(myEndpointJson)))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	cowboyResponse, err := parseRegisterBody(string(resBody))
	if err != nil {
		return err
	}

	g.GameState.SetId(cowboyResponse.Id)
	g.GameState.SetCowboy(*cowboyResponse.Cowboy)

	return nil
}

func parseRegisterBody(body string) (*cowboys.CowboyResponse, error) {
	var cowboyResponse = &cowboys.CowboyResponse{}

	if err := json.Unmarshal([]byte(body), cowboyResponse); err != nil {
		return nil, err
	}

	return cowboyResponse, nil
}

func getMyEndpoint() cowboys.RegisterCowboy {
	host := os.Getenv("COWBOYS_SERVICE_HOST")
	portStr := os.Getenv("COWBOYS_SERVICE_PORT")

	port, err := strconv.Atoi(portStr)

	if err != nil {
		panic("port not a number? how to live what to do (sad)")
	}

	return cowboys.RegisterCowboy{Host: host, Port: port}
}

func New(state *state.GameState, masterEndpoint string) *Game {
	return &Game{GameState: state, GameMasterEdnpoint: masterEndpoint}
}
