package main

import "github.com/gilgalad195/pokedexcli/internal/pokeapi"

func GetLocationUrl(myConfig *pokeapi.Config) string {
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

var WorldMap = map[string][]string{
	"littleroot-town-area": {"hoenn-route-101-area"},
	"hoenn-route-101-area": {"littleroot-town-area", "oldale-town-area"},
	"oldale-town-area":     {"hoenn-route-101-area", "hoenn_route-102-area", "hoenn-route-103-area"},
	"hoenn-route-103-area": {"oldale-town-area"},
	"hoenn-route-102-area": {"petalburg-town-area", "oldale-town-area"},
}
