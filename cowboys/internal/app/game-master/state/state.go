package state

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"cowboys/internal/pkg/cowboys"
)

const (
	Register   = "Register"
	Ready      = "Ready"
	InProgress = "In Progress"
	Done       = "Done"
)

type GameState struct {
	mu                sync.Mutex
	RegisteredPlayers []*cowboys.GameCowboy `json:"registeredPlayers" xml:"registeredPlayers" form:"registeredPlayers"`
	InitialPlayers    []*cowboys.Cowboy     `json:"initialPlayers" xml:"initialPlayers" form:"initialPlayers"`
	PlayersNumbers    int                   `json:"playersNumber" xml:"playersNumber" form:"playersNumber"`
	Status            string                `json:"status" xml:"status" form:"status"`
	cowboysMap        map[string]*cowboys.GameCowboy
}

func (s *GameState) RegisterCowboy(registerData *cowboys.RegisterCowboy) *cowboys.CowboyResponse {
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
		// TODO :D
		panic(err)
	}

	fmt.Printf("Cowboy notified: %s", cowboy.String())
}

func generateId(cowboy *cowboys.Cowboy, registerData *cowboys.RegisterCowboy) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s-%s-%d", cowboy.Name, registerData.Host, registerData.Port)))

	return hex.EncodeToString(hash[:])
}

func New() *GameState {
	// TODO get players from redis
	var initialPlayers []*cowboys.Cowboy
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Eliot", Health: 10, Damage: 1})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Deth", Health: 5, Damage: 5})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Pawl", Health: 5, Damage: 1})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Dvil", Health: 8, Damage: 3})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Gatt", Health: 6, Damage: 1})
	initialPlayers = append(initialPlayers, &cowboys.Cowboy{Name: "Luci", Health: 12, Damage: 2})

	playerNumbers := len(initialPlayers)

	return &GameState{InitialPlayers: initialPlayers,
		PlayersNumbers: playerNumbers, Status: Register,
		cowboysMap: map[string]*cowboys.GameCowboy{}}
}
