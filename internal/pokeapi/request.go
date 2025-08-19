package pokeapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

type Config struct {
	GameVersion      string                 `json:"game_version"`
	CurrentLocation  string                 `json:"current_location"`
	MapCache         *pokecache.Cache       `json:"-"`
	LocationCache    *pokecache.Cache       `json:"-"`
	LastFoundPokemon string                 `json:"-"`
	CaughtPokemon    map[string]PokemonData `json:"caught"`
}

func FetchData(url string, config *Config) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return body, nil
}

func FetchHeaders(url string) (*http.Response, error) {
	req, _ := http.NewRequest("HEAD", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("an error occurred: %v", err)
	}
	return resp, nil
}
