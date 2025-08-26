package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

func main() {

	mapCache := pokecache.NewCache(30 * time.Minute)
	locationCache := pokecache.NewCache(30 * time.Minute)

	myConfig := &gamedata.Config{
		GameVersion:      "sapphire",
		CurrentLocation:  "littleroot-town-area",
		MapCache:         mapCache,
		LocationCache:    locationCache,
		LastFoundPokemon: "",
		CaughtPokemon:    map[string]gamedata.PokemonData{},
		PartyPokemon:     map[int]gamedata.PokemonStatus{},
	}

	//this is my REPL loop, which looks for user input and executes commands
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		hasToken := scanner.Scan()
		if hasToken {
			lowerString := cleanInput(scanner.Text())
			if len(lowerString) == 0 {
				fmt.Println("Enter help if unsure of commands")
			} else {
				cmdInput := lowerString[0]
				args := lowerString[1:]
				if command, exists := commands[cmdInput]; exists {
					err := command.callback(myConfig, args)
					if err != nil {
						fmt.Println("Error:", err)
					}
				} else {
					fmt.Println("Unknown command")
				}
			}
		}
	}
}

func cleanInput(text string) []string {
	cleanText := strings.Fields(strings.ToLower(text))
	return cleanText
}
