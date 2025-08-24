package gamedata

import (
	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

type Config struct {
	GameVersion      string                   `json:"game_version"`
	CurrentLocation  string                   `json:"current_location"`
	MapCache         *pokecache.Cache         `json:"-"`
	LocationCache    *pokecache.Cache         `json:"-"`
	LastFoundPokemon string                   `json:"-"`
	CaughtPokemon    map[string]PokemonData   `json:"caught"`
	PartyPokemon     map[string]PokemonStatus `json:"party_pokemon"`
}
