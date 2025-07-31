package pokeapi

import (
	"encoding/json"
)

type LocationsResp struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func FormatResponse(body []byte) (*LocationsResp, error) {
	var locations LocationsResp
	if err := json.Unmarshal(body, &locations); err != nil {
		return nil, err
	}
	return &locations, nil
}
