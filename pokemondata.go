package main

import (
	"fmt"

	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
)

func CheckValidVersion(versions []string, target string) bool {
	for _, v := range versions {
		if target == v {
			return true
		}
	}
	return false
}

func GetPokemonData(name string, myConfig *pokeapi.Config) (*pokeapi.PokemonData, error) {
	pokeEndpoint := "https://pokeapi.co/api/v2/pokemon/"
	pokeUrl := pokeEndpoint + name

	body, err := pokeapi.FetchData(pokeUrl, myConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}

	pokemon, err := pokeapi.FormatPokemonData(body)
	if err != nil {
		return nil, fmt.Errorf("failed to format response: %v", err)
	}
	return pokemon, nil
}
