package pokeapi

import (
	"encoding/json"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
)

func FormatMapResponse(body []byte) (*gamedata.MapResp, error) {
	var response gamedata.MapResp
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func FormatLocationData(body []byte) (*gamedata.LocationData, error) {
	var response gamedata.LocationData
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func FormatPokemonData(body []byte) (*gamedata.PokemonData, error) {
	var response gamedata.PokemonData
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
