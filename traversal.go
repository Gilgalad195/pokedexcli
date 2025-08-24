package main

import (
	"fmt"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
)

func GetLocationUrl(myConfig *gamedata.Config) string {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	return baseUrl + myConfig.CurrentLocation + "/"
}

func CheckValidPath(paths []string, target string) bool {
	for _, path := range paths {
		if target == path {
			return true
		}
	}
	return false
}

func HasEncounters(url string) (bool, error) {
	resp, err := pokeapi.FetchHeaders(url)
	if err != nil {
		return false, fmt.Errorf("failed to fetch headers: %v", err)
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}

var WorldMap = map[string][]string{
	"littleroot-town-area": {"hoenn-route-101-area"},
	"hoenn-route-101-area": {"littleroot-town-area", "oldale-town-area"},
	"oldale-town-area":     {"hoenn-route-101-area", "hoenn_route-102-area", "hoenn-route-103-area"},
	"hoenn-route-103-area": {"oldale-town-area"},
	"hoenn-route-102-area": {"petalburg-town-area", "oldale-town-area"},
}
