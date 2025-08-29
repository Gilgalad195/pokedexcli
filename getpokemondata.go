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

func GetSummary(myConfig *gamedata.Config, name string) (*gamedata.PokemonStatus, error) {
	var data gamedata.PokemonData
	if caught, ok := myConfig.CaughtPokemon[name]; ok {
		data = caught
	} else {
		fetched, err := GetPokemonData(name, myConfig)
		if err != nil {
			return &gamedata.PokemonStatus{}, fmt.Errorf("unable to get pokemon data: %v", err)
		}
		data = *fetched
	}

	statMap := make(map[string]int)
	for _, stat := range data.Stats {
		statMap[stat.Stat.Name] = stat.BaseStat
	}

	primaryType := data.Types[0].Type.Name

	summary := gamedata.PokemonStatus{
		Name:      data.Name,
		Stats:     statMap,
		CurrentHP: statMap["hp"],
		Fainted:   false,
		Type:      primaryType,
	}

	return &summary, nil
}
