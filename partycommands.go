package main

import (
	"fmt"
	"strconv"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
)

func PartyAdd(myConfig *gamedata.Config, pokename string) {
	if _, exists := myConfig.CaughtPokemon[pokename]; !exists {
		fmt.Printf("You have not caught %s\n", pokename)
		return
	}

	for _, member := range myConfig.PartyPokemon {
		if member.Name == pokename {
			fmt.Printf("%s is already in your party!\n", member.Name)
			return
		}
	}

	summary, err := GetSummary(myConfig, pokename)
	if err != nil {
		fmt.Printf("Failed to get summary: %v", err)
		return
	}

	for i := 1; i <= 6; i++ {
		slotKey := i
		if _, exists := myConfig.PartyPokemon[slotKey]; !exists {
			myConfig.PartyPokemon[slotKey] = summary
			fmt.Printf("Added %s to party slot %d!\n", summary.Name, slotKey)
			return
		}
	}
	fmt.Println("Your party is full. Please remove a pokemon before adding.")

}

func PartyRemove(party map[int]gamedata.PokemonStatus, pokename string) {
	for i, member := range party {
		if member.Name == pokename {
			delete(party, i)
			fmt.Printf("Removed %s from the party\n", pokename)
			return
		}
	}
	fmt.Printf("%s is not in your party!\n", pokename)
}

func PartyInspect(party map[int]gamedata.PokemonStatus, pokename string) {
	for i, member := range party {
		if member.Name == pokename {
			fmt.Printf("Slot %d:\n", i)
			fmt.Printf(" - Name: %s\n", member.Name)
			fmt.Printf(" - HP: %d/%d\n", member.CurrentHP, member.Stats["hp"])
			fmt.Printf(" - Attack: %d\n", member.Stats["attack"])
			fmt.Printf(" - Defense: %d\n", member.Stats["defense"])
			fmt.Printf(" - Sp Attack: %d\n", member.Stats["special-attack"])
			fmt.Printf(" - Sp Defense: %d\n", member.Stats["special-defense"])
			fmt.Printf(" - Speed: %d\n", member.Stats["speed"])
			return
		}
	}
	fmt.Printf("%s is not in your party!\n", pokename)
}

func PartySwap(party map[int]gamedata.PokemonStatus, a, b string) {
	intA, errA := strconv.Atoi(a)
	intB, errB := strconv.Atoi(b)
	if errA != nil || errB != nil || intA < 1 || intA > 6 || intB < 1 || intB > 6 {
		fmt.Println("Invalid slot numbers. Please use numeric values of 1 - 6")
		return
	}
	fmt.Printf("Swapping slots %s and %s\n", a, b)
	party[intA], party[intB] = party[intB], party[intA]
}

func PartyList(party map[int]gamedata.PokemonStatus) {
	for i := 1; i <= 6; i++ {
		member := party[i]
		fmt.Printf("Slot %d: %s\n", i, member.Name)
	}
}
