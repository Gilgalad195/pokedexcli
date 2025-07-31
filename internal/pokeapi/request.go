package pokeapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

type Config struct {
	Next          string
	Previous      string
	MapCache      *pokecache.Cache
	LocationCache *pokecache.Cache
	Search        string
}

func GetLocations(url string, config *Config) ([]byte, error) {
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
