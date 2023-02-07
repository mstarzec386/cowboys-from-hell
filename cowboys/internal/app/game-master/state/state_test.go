package state

import (
	"errors"
	"fmt"
	"testing"
)

type RedisMockedClient struct {
	data map[string]string
}

func (r *RedisMockedClient) Get(key string) (string, error) {
	if data := r.data[key]; data != "" {
		return data, nil
	}

	return "", errors.New("not found")
}

func (r *RedisMockedClient) Set(key string, data interface{}) error {
	if r.data == nil {
		r.data = make(map[string]string)
	}

	r.data[key] = fmt.Sprint(data)

	return nil
}

func TestNew(t *testing.T) {
	t.Run("Fallback", func(t *testing.T) {
		redisMockedClient := RedisMockedClient{}
		state := New(&redisMockedClient)

		t.Logf("Initial players %v", state.InitialPlayers)

		if numberOfInitialPlayers := len(state.InitialPlayers); numberOfInitialPlayers != 5 {
			t.Errorf("Number of players incorrect, got: %d, want: %d.", numberOfInitialPlayers, 5)
		}

		if firstCowboy := state.InitialPlayers[0]; firstCowboy.Name != "Elto" {
			t.Errorf("Wrong name of first player, got: %s, want: %s.", firstCowboy.Name, "Elto")
		}
	})

	t.Run("From redis", func(t *testing.T) {
		redisMockedClient := RedisMockedClient{}
		redisMockedClient.Set("init", "[{\"name\": \"Weirdo\", \"health\": 667, \"damage\": 13}]")

		state := New(&redisMockedClient)

		t.Logf("Initial players %v", state.InitialPlayers)

		if numberOfInitialPlayers := len(state.InitialPlayers); numberOfInitialPlayers != 1 {
			t.Errorf("Number of players incorrect, got: %d, want: %d.", numberOfInitialPlayers, 1)
		}

		firstCowboy := state.InitialPlayers[0]

		if firstCowboy.Name != "Weirdo" {
			t.Errorf("Wrong name of first player, got: %s, want: %s.", firstCowboy.Name, "Weirdo")
		}
		if firstCowboy.Health != 667 {
			t.Errorf("Wrong health of first player, got: %d, want: %d.", firstCowboy.Health, 667)
		}
		if firstCowboy.Damage != 13 {
			t.Errorf("Wrong damage of first player, got: %d, want: %d.", firstCowboy.Damage, 13)
		}
	})
}
