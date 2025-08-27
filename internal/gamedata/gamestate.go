package gamedata

import (
	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

type Config struct {
	GameVersion        string                 `json:"game_version"`
	CurrentLocation    string                 `json:"current_location"`
	MapCache           *pokecache.Cache       `json:"-"`
	LocationCache      *pokecache.Cache       `json:"-"`
	EncounteredPokemon PokemonStatus          `json:"-"`
	CaughtPokemon      map[string]PokemonData `json:"caught"`
	PartyPokemon       map[int]PokemonStatus  `json:"party_pokemon"`
}
