package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
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

func commandExit(_ *gamedata.Config, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *gamedata.Config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(myConfig *gamedata.Config, _ []string) error {
	fmt.Println("Your map is concerningly empty.")
	return nil
}

func commandLook(myConfig *gamedata.Config, _ []string) error {
	fmt.Printf("You look around '%v' and see the following paths:\n", myConfig.CurrentLocation)
	availablePaths := WorldMap[myConfig.CurrentLocation]
	for _, path := range availablePaths {
		fmt.Printf(" - %s\n", path)
	}
	return nil
}

func commandMove(myConfig *gamedata.Config, args []string) error {
	if !CheckValidPath(WorldMap[myConfig.CurrentLocation], args[0]) {
		fmt.Println("You can't get there from here!")
		return nil
	} else {
		myConfig.CurrentLocation = args[0]
		fmt.Printf("You are now in '%s'\n", myConfig.CurrentLocation)
		myConfig.LastFoundPokemon = ""
	}
	return nil
}

func commandExplore(myConfig *gamedata.Config, _ []string) error {
	locationUrl := GetLocationUrl(myConfig)

	hasEncounters, err := HasEncounters(locationUrl)
	if err != nil {
		fmt.Printf("unable to determine encounter status: %v\n", err)
	} else if !hasEncounters {
		fmt.Println("there are no wild Pokémon to find here")
	} else {

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
		var availablePokemon []string

		fmt.Println("The following pokemon can be found here:")
		for _, encounter := range encounters {
			var versions []string
			for _, vd := range encounter.VersionDetails {
				versions = append(versions, vd.Version.Name)
			}

			if CheckValidVersion(versions, myConfig.GameVersion) {
				fmt.Printf(" - %s\n", encounter.Pokemon.Name)
				availablePokemon = append(availablePokemon, encounter.Pokemon.Name)
			}
		}

		if len(availablePokemon) > 0 {
			roll := rand.Intn(len(availablePokemon))
			foundPokemon := availablePokemon[roll]
			myConfig.LastFoundPokemon = foundPokemon
			fmt.Println("You found a wild Pokemon!")
			fmt.Printf(" - %s\n\n", foundPokemon)
		} else {
			fmt.Println(" - there are no pokemon to find")
		}
	}
	return nil
}

func commandCatch(myConfig *gamedata.Config, args []string) error {
	if len(args) == 0 {
		fmt.Println("please enter a pokemon name")
	} else {
		pokeName := args[0]
		if pokeName != myConfig.LastFoundPokemon {
			fmt.Println("that pokemon isn't here right now")
		} else {
			pokemon, err := GetPokemonData(pokeName, myConfig)
			if err != nil {
				return fmt.Errorf("failed to get pokemon data: %v", err)
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
	}

	return nil
}

func commandInspect(myConfig *gamedata.Config, args []string) error {
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

func commandPokedex(myConfig *gamedata.Config, _ []string) error {
	fmt.Println("Your Pokedex:")
	for name := range myConfig.CaughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandSave(myConfig *gamedata.Config, _ []string) error {
	jsonSave, err := json.Marshal(myConfig)
	if err != nil {
		return fmt.Errorf("save failed: %v", err)
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("save failed. config dir error: %v", err)
	}

	saveDir := filepath.Join(configDir, "PokedexCLI", "saves")
	saveFilePath := filepath.Join(saveDir, "pokesave.json")
	err = os.MkdirAll(saveDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating save dir: %v", err)
	}
	err = os.WriteFile(saveFilePath, jsonSave, 0644)
	if err != nil {
		return fmt.Errorf("error writing to save: %v", err)
	}
	fmt.Println("Save successful!")

	return nil
}

func commandLoad(myConfig *gamedata.Config, _ []string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("config dir error: %v", err)
	}
	saveDir := filepath.Join(configDir, "PokedexCLI", "saves")
	saveFilePath := filepath.Join(saveDir, "pokesave.json")

	data, err := os.ReadFile(saveFilePath)
	if err != nil {
		return fmt.Errorf("save read failed: %v", err)
	}

	err = json.Unmarshal(data, &myConfig)
	if err != nil {
		return fmt.Errorf("load state failed: %v", err)
	}
	fmt.Println("Save data loaded!")

	return nil
}

func commandParty(myConfig *gamedata.Config, args []string) error {
	subcommand := args[0]

	if subcommand == "" || len(args) < 1 {
		fmt.Println("Valid arguments:")
		fmt.Println(" - add/remove/inspect followed by pokemon name")
		fmt.Println(" - swap followed by two party slots")
		fmt.Println(" - list")
		return nil
	}

	if subcommand == "add" {
		if len(args) < 2 || args[1] == "" {
			fmt.Println("Please indicate the caught pokemon you wish to add to your party")
			return nil
		}
		pokename := args[1]
		PartyAdd(myConfig, pokename)
	}

	if subcommand == "remove" {

		if len(args) < 2 || args[1] == "" {
			fmt.Println("Please indicate the pokemon you wish to remove from your party")
			return nil
		}
		pokename := args[1]
		PartyRemove(myConfig.PartyPokemon, pokename)
	}

	if subcommand == "inspect" {
		if len(args) < 2 || args[1] == "" {
			fmt.Println("Please indicate the party pokemon you wish to inspect")
			return nil
		}
		pokename := args[1]
		PartyInspect(myConfig.PartyPokemon, pokename)
	}

	if subcommand == "swap" {
		if len(args) < 3 {
			fmt.Println("Swap requires two slot numbers (1-6)")
			return nil
		}
		a, b := args[1], args[2]
		PartySwap(myConfig.PartyPokemon, a, b)
	}

	if subcommand == "list" {
		PartyList(myConfig.PartyPokemon)
	}
	return nil
}
