package pokeapi

import (
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

func GetLocations(url string, config *Config) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
