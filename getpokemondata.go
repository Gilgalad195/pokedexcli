package main

import (
	"fmt"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
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

func GetPokemonData(name string, myConfig *gamedata.Config) (*gamedata.PokemonData, error) {
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

func GetSummary(myConfig *gamedata.Config, name string) gamedata.PokemonStatus {
	statMap := make(map[string]int)
	for _, stat := range myConfig.CaughtPokemon[name].Stats {
		statMap[stat.Stat.Name] = stat.BaseStat
	}

	summary := gamedata.PokemonStatus{
		Name:      myConfig.CaughtPokemon[name].Name,
		Stats:     statMap,
		CurrentHP: statMap["hp"],
		Fainted:   false,
	}

	return summary
}
