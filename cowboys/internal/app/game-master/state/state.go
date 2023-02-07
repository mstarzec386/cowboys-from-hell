package state

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cowboys/internal/pkg/clients/redis"
	"cowboys/internal/pkg/cowboys"
)

const (
	Register   = "Register"
	Ready      = "Ready"
	InProgress = "In Progress"
	Done       = "Done"
)

type GameState struct {
	RegisteredPlayers []*cowboys.GameCowboy `json:"registeredPlayers" xml:"registeredPlayers" form:"registeredPlayers"`
	InitialPlayers    []*cowboys.Cowboy     `json:"initialPlayers" xml:"initialPlayers" form:"initialPlayers"`
	PlayersNumbers    int                   `json:"playersNumber" xml:"playersNumber" form:"playersNumber"`
	Status            string                `json:"status" xml:"status" form:"status"`
	mu                sync.Mutex
	cowboysMap        map[string]*cowboys.GameCowboy
}

func (s *GameState) RegisterCowboy(registerData *cowboys.RegisterCowboy) *cowboys.CowboyResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	lastOneIndex := len(s.InitialPlayers)

	if lastOneIndex > 0 {
		newCowboy := s.InitialPlayers[lastOneIndex-1]
		s.InitialPlayers = s.InitialPlayers[:lastOneIndex-1]

		gameCowboy := cowboys.GameCowboy{Cowboy: newCowboy, Endpoint: registerData, Id: generateId(newCowboy, registerData)}
		s.RegisteredPlayers = append(s.RegisteredPlayers, &gameCowboy)

		s.cowboysMap[gameCowboy.Id] = &gameCowboy

		if len(s.RegisteredPlayers) == s.PlayersNumbers {
			s.setInprogressStatus()
			go s.notifyCowboys()
		}

		return &cowboys.CowboyResponse{Id: gameCowboy.Id, Cowboy: newCowboy}
	}

	return nil
}

func (s *GameState) UpdateCowboy(id string, updateData *cowboys.UpdateCowboy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO add if state == register then should not allow
	cowboy := s.cowboysMap[id]
	if cowboy == nil {
		return errors.New("not found")
	}

	cowboy.Cowboy.Health = updateData.Health

	if s.allDead() {
		fmt.Println("No winner this time, they are all dead, poor guys ―(x_x)→")
	}

	return nil
}

func (s *GameState) allDead() bool {
	allDead := true

	for _, cowboy := range s.RegisteredPlayers {
		if cowboy.Cowboy.Health > 0 {
			allDead = false
		}

	}

	return allDead
}

func (s *GameState) setInprogressStatus() {
	s.setStatus(Ready)
}

func (s *GameState) notifyCowboys() {
	for _, cowboy := range s.RegisteredPlayers {
		// TODO error handling wait for all responses etc
		go notifyCowboy(cowboy)
	}
}

func (s *GameState) setStatus(status string) {
	s.Status = status
}

func notifyCowboy(cowboy *cowboys.GameCowboy) {
	cowboyUrl := cowboy.Endpoint.ToUrl("start")
	resp, err := http.Get(cowboyUrl)
	if err != nil || resp.StatusCode != 200 {
		// TODO what to do then?
		log.Printf("Notify cowboy failed %s\n", err.Error())
	}

	fmt.Printf("Cowboy notified: %s\n", cowboy.String())
}

func generateId(cowboy *cowboys.Cowboy, registerData *cowboys.RegisterCowboy) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s-%s-%d", cowboy.Name, registerData.Host, registerData.Port)))

	return hex.EncodeToString(hash[:])
}

func New(redisClient redis.RedisClientInterface) *GameState {
	initialPlayers, err := getInitialPlayers(redisClient)

	// FALLBACK if redis is not working
	if err != nil {
		fmt.Printf("Can't get initial data form redis, static cowboys loaded: %s\n", err.Error())

		initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Elto", Health: 10, Damage: 1})
		initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Deth", Health: 5, Damage: 5})
		initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Pawl", Health: 5, Damage: 1})
		initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Dvil", Health: 8, Damage: 3})
		initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Luci", Health: 12, Damage: 2})

	}

	playerNumbers := len(initialPlayers)

	return &GameState{InitialPlayers: initialPlayers,
		PlayersNumbers: playerNumbers, Status: Register,
		cowboysMap: map[string]*cowboys.GameCowboy{}}
}

func getInitialPlayers(redisClient redis.RedisClientInterface) ([]*cowboys.Cowboy, error) {
	var initialPlayers = []*cowboys.Cowboy{}
	data, err := redisClient.Get("init")
	if err != nil {
		return initialPlayers, err
	}

	if err := json.Unmarshal([]byte(data), &initialPlayers); err != nil {
		return initialPlayers, err
	}

	return initialPlayers, nil
}
