package game

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"cowboys/internal/app/cowboy/state"
	"cowboys/internal/pkg/cowboys"
)

type Game struct {
	GameMasterEdnpoint string
	GameState          *state.GameState
	Id                 string
}

func (g *Game) GetGameState() *state.GameState{
	return g.GameState
}

func (g *Game) Register() error {
	fmt.Println("register")

	registerEndpoint := fmt.Sprintf("%s/cowboys", g.GameMasterEdnpoint)
	myEndpoint := getMyEndpoint()
	myEndpointJson, err := json.Marshal(myEndpoint)
	if err != nil {
		return err
	}

	fmt.Printf("register: %s with %s",registerEndpoint, myEndpointJson)

	res, err := http.Post(registerEndpoint, "application/json", bytes.NewBuffer([]byte(myEndpointJson)))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)
	cowboyResponse, err := parseRegisterBody(string(resBody))
	if err != nil {
		return err
	}
	fmt.Println(*cowboyResponse)


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
