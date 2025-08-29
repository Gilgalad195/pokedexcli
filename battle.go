package main

import (
	"fmt"
	"math/rand"
	"slices"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
)

var superEffectiveMap = map[string][]string{
	"normal":   {},
	"fire":     {"grass", "ice", "bug", "steel"},
	"water":    {"fire", "ground", "rock"},
	"electric": {"water", "flying"},
	"grass":    {"water", "ground", "rock"},
	"ice":      {"grass", "ground", "flying", "dragon"},
	"fighting": {"normal", "ice", "rock", "steel"},
	"poison":   {"poison"},
	"ground":   {"ground", "electric", "poison", "rock", "steel"},
	"flying":   {"grass", "fighting", "bug"},
	"psychic":  {"psychic", "poison"},
	"bug":      {"grass", "psychic", "dark"},
	"rock":     {"fire", "grass", "ground", "bug"},
	"ghost":    {"psychic", "ghost"},
	"dragon":   {"dragon"},
	"dark":     {"psychic", "ghost"},
	"steel":    {"ice", "rock"},
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// effectiveness is going to use TCG logic, only counting super effective matchups
func effectiveness(attacking *gamedata.PokemonStatus, defending *gamedata.PokemonStatus) int {
	if slices.Contains(superEffectiveMap[attacking.Type], defending.Type) {
		return 2
	}
	return 1
}

func TakeDamage(defending *gamedata.PokemonStatus, dam int) {
	defending.CurrentHP -= dam
	if defending.CurrentHP <= 0 {
		defending.Fainted = true
		defending.CurrentHP = 0
	}
}

func PokemonAttack(attacking *gamedata.PokemonStatus, defending *gamedata.PokemonStatus) {
	level := 1
	crit := 1
	if rand.Intn(16) == 0 { // ~6.25% chance
		crit = 2
	}
	power := 60
	var a int
	var d int
	if attacking.Stats["attack"] >= attacking.Stats["special-attack"] {
		a = attacking.Stats["attack"]
		d = defending.Stats["defense"]
	} else {
		a = attacking.Stats["special-attack"]
		d = defending.Stats["special-defense"]
	}
	stab := 1
	type1 := effectiveness(attacking, defending)
	type2 := 1
	random := float64(RandomInt(217, 255)) / 255.0

	base := (((2.0*float64(level)*float64(crit))/5.0)+2.0)*float64(power)*(float64(a)/float64(d))/50.0 + 2.0
	damage := int(base * float64(stab) * float64(type1) * float64(type2) * random)

	TakeDamage(defending, int(damage))

	fmt.Printf("%s hit %s for %v damage!\n", attacking.Name, defending.Name, int(damage))
	fmt.Printf("%s has %d HP remaining.\n", defending.Name, defending.CurrentHP)
	if defending.Fainted {
		fmt.Printf("%s has fainted!\n", defending.Name)
	}
}
