package main

import (
	"fmt"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
)

type LocationArea struct {
	North string
	East  string
	South string
	West  string
}

func GetLocationUrl(myConfig *gamedata.Config) string {
	baseUrl := "https://pokeapi.co/api/v2/location-area/"
	return baseUrl + myConfig.CurrentLocation + "/"
}

func GetDirections(here LocationArea, there string) string {
	switch there {
	case "north":
		return here.North
	case "east":
		return here.East
	case "south":
		return here.South
	case "west":
		return here.West
	default:
		return ""
	}
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
	"oldale-town-area":     {"hoenn-route-101-area", "hoenn-route-102-area", "hoenn-route-103-area"},
	"hoenn-route-103-area": {"oldale-town-area"},
	"hoenn-route-102-area": {"petalburg-town-area", "oldale-town-area"},
}

var WorldMapV2 = map[string]LocationArea{
	"littleroot-town-area": {
		North: "hoenn-route-101-area",
	},
	"hoenn-route-101-area": {
		North: "oldale-town-area",
		South: "littleroot-town-area",
	},
	"oldale-town-area": {
		North: "hoenn-route-103-area",
		South: "hoenn-route-101-area",
		West:  "hoenn-route-102-area",
	},

	// Route 102 → Petalburg City → Route 104 South
	"hoenn-route-102-area": {
		West: "petalburg-city-area",
		East: "oldale-town-area",
	},
	"hoenn-route-103-area": {
		South: "Oldale Town",
	},
	"petalburg-city-area": {
		West: "hoenn-route-104-area",
		East: "hoenn-route-102-area",
	},
	"hoenn-route-104-area": {
		North: "rustboro-city-area",
		South: "petalburg-city-area",
		West:  "petalburg-woods-area",
	},
	"petalburg-woods-area": {
		East: "hoenn-route-104-area",
	},
	"rustboro-city-area": {
		North: "hoenn-route-115-area",
		South: "hoenn-route-104-north-area",
		East:  "hoenn-route-116-area",
	},

	"hoenn-route-115-area": {
		South: "rustboro-city-area",
	},

	// Route 116 → Rusturf Tunnel → Verdanturf Town → Route 117
	"hoenn-route-116-area": {
		East: "rusturf-tunnel-area",
		West: "rustboro-city-area",
	},
	"rusturf-tunnel-area": {
		South: "verdanturf-town-area",
		West:  "hoenn-route-116-area",
	},
	"verdanturf-town-area": {
		North: "rusturf-tunnel-area",
		East:  "hoenn-route-117-area",
	},
	"hoenn-route-117-area": {
		East: "mauville-city-area",
		West: "verdanturf-town-area",
	},

	// Mauville City between Route 110 and Route 111
	"mauville-city-area": {
		North: "hoenn-route-111-area",
		South: "hoenn-route-110-area",
		West:  "hoenn-route-117-area",
		// East intentionally omitted (Route 118)
	},
	"hoenn-route-110-area": {
		North: "mauville-city-area",
		South: "slateport-city-area",
	},
	"slateport-city-area": {
		North: "hoenn-route-110-area",
		South: "hoenn-route-109-area",
	},
	"hoenn-route-111-area": {
		North: "hoenn-route-113-area",
		South: "mauville-city-area",
		West:  "hoenn-route-112-area",
	},

	// Route 112 → Fiery Path → Route 113 → Fallarbor Town → Meteor Falls
	"hoenn-route-112-area": {
		North: "fiery-path-area",
		South: "hoenn-route-111-area",
		West:  "jagged-pass-area",
	},
	"jagged-pass-area": {
		North: "mt-chimney-area",
		West:  "lavaridge-town-area",
	},
	"mount-chimney-area": {
		South: "jagged-pass-area",
	},
	"lavaridge-town-area": {
		East: "jagged-pass-area",
	},
	"fiery-path-area": {
		North: "hoenn-route-113-area",
		South: "hoenn-route-112-area",
	},
	"hoenn-route-113-area": {
		North: "fallarbor-town-area",
		South: "fiery-path-area",
		East:  "hoenn-route-111-area",
	},
	"fallarbor-town-area": {
		East: "hoenn-route-113-area",
		West: "hoenn-route-114-area",
	},
	"hoenn-route-114-area": {
		East: "fallarbor-town-area",
		West: "meteor-falls-area",
	},
	"meteor-falls-area": {
		East:  "hoenn-route-114-area",
		South: "hoenn-route-115-area",
	},

	// Route 106 → Dewford Town → Route 107 → Route 108 → Slateport
	"hoenn-route-106-area": {
		East:  "dewford-town-area",
		South: "granite-cave-area",
	},
	"granite-cave-area": {
		North: "hoenn-route-106-area",
	},
	"dewford-town-area": {
		North: "hoenn-route-107-area",
		West:  "hoenn-route-106-area",
	},
	"hoenn-route-107-area": {
		East: "hoenn-route-108-area",
		West: "dewford-town-area",
	},
	"hoenn-route-108-area": {
		East: "hoenn-route-109-area",
		West: "hoenn-route-107-area",
	},
}
