package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

func main() {

	mapCache := pokecache.NewCache(30 * time.Minute)
	locationCache := pokecache.NewCache(30 * time.Minute)

	myConfig := &pokeapi.Config{
		//this initializes the struct with the starting API url.
		BaseUrl:       "https://pokeapi.co/api/v2/location-area/",
		Next:          "https://pokeapi.co/api/v2/location-area/",
		Previous:      "",
		MapCache:      mapCache,
		LocationCache: locationCache,
		CaughtPokemon: map[string]pokeapi.PokemonData{},
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

func commandExit(_ *pokeapi.Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *pokeapi.Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(myConfig *pokeapi.Config, _ []string) error {
	if myConfig.Next == "" {
		fmt.Println("you're on the last page")
	} else {
		//declaring body and err for use later
		var body []byte
		var err error
		val, found := myConfig.MapCache.Get(myConfig.Next)
		if found {
			fmt.Println("This was in cache!")
			body = val
		} else {
			body, err = pokeapi.FetchData(myConfig.Next, myConfig)
			if err != nil {
				return fmt.Errorf("failed to fetch data: %v", err)
			}
			myConfig.MapCache.Add(myConfig.Next, body)
		}

		data, err := pokeapi.FormatMapResponse(body)
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}
		locations := data.Results
		for i, location := range locations {
			if i >= 20 {
				break
			}
			fmt.Println(location.Name)
		}
		myConfig.Previous = myConfig.Next
		myConfig.Next = data.Next
	}
	return nil
}

func commandMapb(myConfig *pokeapi.Config, _ []string) error {
	if myConfig.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		var body []byte
		var err error
		val, found := myConfig.MapCache.Get(myConfig.Previous)
		if found {
			fmt.Println("This was in cache!")
			body = val
		} else {
			body, err = pokeapi.FetchData(myConfig.Previous, myConfig)
			if err != nil {
				return fmt.Errorf("failed to fetch data: %v", err)
			}
			myConfig.MapCache.Add(myConfig.Previous, body)
		}

		data, err := pokeapi.FormatMapResponse(body)
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}
		locations := data.Results
		for i, location := range locations {
			if i >= 20 {
				break
			}
			fmt.Println(location.Name)
		}
		myConfig.Next = myConfig.Previous
		myConfig.Previous = data.Previous
	}
	return nil
}

func commandExplore(myConfig *pokeapi.Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please enter a location")
	} else {

		locationUrl := myConfig.BaseUrl + args[0]

		var body []byte
		var err error
		val, found := myConfig.LocationCache.Get(locationUrl)
		if found {
			fmt.Println("This was in cache!")
			body = val
		} else {
			body, err = pokeapi.FetchData(locationUrl, myConfig)
			if err != nil {
				return fmt.Errorf("failed to fetch data: %v", err)
			}
			myConfig.LocationCache.Add(locationUrl, body)
		}

		//data is the formatted LocationData struct
		data, err := pokeapi.FormatLocationData(body)
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}
		encounters := data.PokemonEncounters
		fmt.Println("Found Pokemon:")
		for _, pokemon := range encounters {
			fmt.Printf("- %s\n", pokemon.Pokemon.Name)
		}
	}
	return nil
}

func commandCatch(myConfig *pokeapi.Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please enter a pokemon name")
	} else {
		pokeEndpoint := "https://pokeapi.co/api/v2/pokemon/"
		pokeName := args[0]
		pokeUrl := pokeEndpoint + pokeName

		body, err := pokeapi.FetchData(pokeUrl, myConfig)
		if err != nil {
			return fmt.Errorf("failed to fetch data: %v", err)
		}

		pokemon, err := pokeapi.FormatPokemonData(body)
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)
		roll := rand.Intn(pokemon.BaseExperience)
		var success bool
		if roll <= 40 {
			success = true
		}

		if success {
			fmt.Printf("%s was caught!\n", pokeName)
			myConfig.CaughtPokemon[pokeName] = *pokemon
		} else {
			fmt.Printf("%s escaped!\n", pokeName)
		}
	}

	return nil
}

func commandInspect(myConfig *pokeapi.Config, args []string) error {
	pokeName := args[0]
	if pokemon, exists := myConfig.CaughtPokemon[pokeName]; exists {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, poketype := range pokemon.Types {
			fmt.Printf("  - %s\n", poketype.Type.Name)
		}
	} else {
		fmt.Printf("you have not caught that pokemon\n")
	}
	return nil
}

func commandPokedex(myConfig *pokeapi.Config, _ []string) error {
	fmt.Println("Your Pokedex:")
	for name := range myConfig.CaughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandSave(myConfig *pokeapi.Config, _ []string) error {
	jsonSave, err := json.Marshal(myConfig)
	if err != nil {
		return fmt.Errorf("save failed: %v", err)
	}

	fmt.Println(string(jsonSave))
	//this is where I need to add OS manipulation. Research how to do that in Go.
	return nil
}
