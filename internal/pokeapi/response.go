package pokeapi

import (
	"encoding/json"
	"io"
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

func FormatResponse(body io.Reader) (*LocationsResp, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	var locations LocationsResp
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, err
	}
	return &locations, nil
}
