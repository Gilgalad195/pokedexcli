package pokeapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
)

func FetchData(url string, config *gamedata.Config) ([]byte, error) {
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
