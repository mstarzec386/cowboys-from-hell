package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
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
	gameCowboys        []*cowboys.GameCowboy
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

func (g *Game) GetDamage() int {
	return g.GameState.GetDamage()
}

func (g *Game) HitCowboy(damage int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.GameState.HitCowboy(damage)
	g.raportState()

	fmt.Printf("Oh someones hits (%d) me (%d) ಥ_ಥ", damage, g.GetHealth())
	if g.GameState.GetHealth() < 1 {
		fmt.Println(" and killed (✖╭╮✖)")
		g.Stop()
	}
	fmt.Println("")

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

func (g *Game) Stop() error {
	if g.ticker == nil {
		return errors.New("game already stopped")
	}

	g.done <- true
	g.ticker = nil

	return nil
}

func (g *Game) IsRegistered() bool {
	if g.GameState.GetId() != "" {
		return true
	} else {
		return false
	}
}

func (g *Game) Register() error {
	registerEndpoint := fmt.Sprintf("%s/cowboys", g.GameMasterEdnpoint)
	myEndpoint := getMyEndpoint()
	myEndpointJson, err := json.Marshal(myEndpoint)
	if err != nil {
		return err
	}

	// fmt.Printf("register: %s with %s\n", registerEndpoint, myEndpointJson)

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

	fmt.Printf("Registered as %s\n", cowboyResponse.Cowboy.String())
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

func parseCowboysBody(body string) ([]*cowboys.GameCowboy, error) {
	var cowboys []*cowboys.GameCowboy

	if err := json.Unmarshal([]byte(body), &cowboys); err != nil {
		return nil, err
	}

	return cowboys, nil
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

func (g *Game) runLoop() {
	g.mu.Lock()
	defer g.mu.Unlock()

	ticker := time.NewTicker(1 * time.Second)
	g.ticker = ticker

	done := make(chan bool)
	g.done = done

	go func() {
		for {
			select {
			case <-ticker.C:
				g.letsPlay()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *Game) letsPlay() {
	// fmt.Println("Next round")
	fmt.Println("")
	g.pullCowboysState()
	g.shoot()
	// if g.GameState.GetHealth() < 1 {
	// 	g.Stop()
	// }
}

func (g *Game) filterCowboysState(cowboysState []*cowboys.GameCowboy) []*cowboys.GameCowboy {
	var cowboysStateFiltered []*cowboys.GameCowboy

	for _, cowboy := range cowboysState {
		if cowboy.Id != g.GameState.GetId() && cowboy.Cowboy.Health > 0 {
			cowboysStateFiltered = append(cowboysStateFiltered, cowboy)
		}
	}

	return cowboysStateFiltered
}

func (g *Game) pullCowboysState() error {
	registerEndpoint := fmt.Sprintf("%s/cowboys", g.GameMasterEdnpoint)

	res, err := http.Get(registerEndpoint)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, _ := io.ReadAll(res.Body)
	cowboysState, err := parseCowboysBody(string(resBody))
	if err != nil {
		return err
	}

	g.gameCowboys = g.filterCowboysState(cowboysState)

	return nil

}

func (g *Game) shoot() {
	if g.GameState.GetHealth() < 1 {
		return
	}

	numberOfPlayers := len(g.gameCowboys)

	if numberOfPlayers < 1 {
		fmt.Printf("I WIN!!!\nFck yea motherfuckers!\n\\(ᵔᵕᵔ)/   \\(ᵔᵕᵔ)/   \\(ᵔᵕᵔ)/\n")

		g.Stop()
		return
	}

	randomIndex := getRandom(numberOfPlayers)
	victim := g.gameCowboys[randomIndex]

	// would be nice to have method like this
	// victim.Hit()

	hitVictim(victim, g.GameState.GetDamage())
}

func (g *Game) raportState() error {
	registerEndpoint := fmt.Sprintf("%s/cowboys/%s", g.GameMasterEdnpoint, g.GameState.GetId())
	myStateJson, err := json.Marshal(g.GetCowboy())
	if err != nil {
		return err
	}

	// fmt.Printf("Raport state: %s with %s\n", registerEndpoint, myStateJson)
	res, err := PutRequest(registerEndpoint, "application/json", bytes.NewBuffer([]byte(myStateJson)))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// resBody, _ := io.ReadAll(res.Body)

	return nil
}

func PutRequest(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", contentType)
	res, err := client.Do(req)

	return res, err
}

func hitVictim(victim *cowboys.GameCowboy, damage int) error {
	// TODO move paths to const
	fmt.Printf("Hit victim %s (%d)\n", victim.Cowboy.Name, victim.Cowboy.Health)
	res, err := http.Get(victim.Endpoint.ToUrl(fmt.Sprintf("hit/%d", damage)))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// resBody, _ := io.ReadAll(res.Body)

	if res.StatusCode == 400 {
		fmt.Printf("Shit, He is already dead\n")
	}

	// fmt.Println("Response from victim", res.StatusCode, string(resBody))

	return nil
}

// TODO move to helpers
func getRandom(n int) int {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	return random.Intn(n)
}

func New(state *state.GameState, masterEndpoint string) *Game {
	return &Game{GameState: state, GameMasterEdnpoint: masterEndpoint}
}