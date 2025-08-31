package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/gilgalad195/pokedexcli/internal/gamedata"
	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
)

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
	if loc, ok := WorldMapV2[myConfig.CurrentLocation]; ok {
		if loc.North != "" {
			fmt.Printf("North: %s\n", loc.North)
		}
		if loc.East != "" {
			fmt.Printf("East: %s\n", loc.East)
		}
		if loc.South != "" {
			fmt.Printf("South: %s\n", loc.South)
		}
		if loc.West != "" {
			fmt.Printf("West: %s\n", loc.West)
		}
		fmt.Println("")
	} else {
		fmt.Println("Problem getting location")
	}
	return nil
}

func commandMove(myConfig *gamedata.Config, args []string) error {
	if len(args) < 1 {
		fmt.Println("You move in place and wind up right where you began")
		return nil
	}
	here := WorldMapV2[myConfig.CurrentLocation]
	there := GetDirections(here, strings.ToLower(args[0]))
	if there != "" {
		myConfig.CurrentLocation = there
		fmt.Printf("You are now in '%s'\n", myConfig.CurrentLocation)
	} else {
		fmt.Println("You can't go that way!")
	}

	myConfig.EncounteredPokemon = nil
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
			myConfig.EncounteredPokemon, err = GetSummary(myConfig, foundPokemon)
			if err != nil {
				return fmt.Errorf("explore failed: %v", err)
			}
			fmt.Printf("A wild %s appeared!\n", foundPokemon)
			fmt.Printf("Wild %s has %d/%d HP\n\n", myConfig.EncounteredPokemon.Name, myConfig.EncounteredPokemon.CurrentHP, myConfig.EncounteredPokemon.Stats["hp"])
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
		if pokeName != myConfig.EncounteredPokemon.Name || myConfig.EncounteredPokemon == nil {
			fmt.Println("that pokemon isn't here right now")
		} else {
			pokemon, err := GetPokemonData(pokeName, myConfig)
			if err != nil {
				return fmt.Errorf("failed to get pokemon data: %v", err)
			}

			fmt.Printf("Throwing a Pokeball at %s...\n", pokeName)
			roll := rand.Intn(pokemon.BaseExperience + myConfig.EncounteredPokemon.CurrentHP)
			var success bool
			if roll <= 50 {
				success = true
			}

			if success {
				fmt.Printf("%s was caught!\n", pokeName)
				myConfig.CaughtPokemon[pokeName] = *pokemon
				myConfig.EncounteredPokemon = nil
			} else {
				fmt.Printf("%s broke free!\n", pokeName)
				mine := myConfig.PartyPokemon[1]
				theirs := myConfig.EncounteredPokemon
				PokemonAttack(theirs, mine)
				if mine.Fainted {
					fmt.Printf("%s got away!", theirs.Name)
					myConfig.EncounteredPokemon = nil
				}
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
	if len(args) < 1 || args[0] == "" {
		fmt.Println("Valid arguments:")
		fmt.Println(" - add/remove/inspect/heal followed by pokemon name")
		fmt.Println(" - swap followed by two party slots")
		fmt.Println(" - list")
		return nil
	}
	subcommand := args[0]

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

	if subcommand == "heal" {
		pokename := args[1]
		PartyHeal(myConfig, pokename)
	}
	return nil
}

func commandRun(myConfig *gamedata.Config, _ []string) error {
	if myConfig.EncounteredPokemon == nil {
		fmt.Println("You are not in an encounter.")
	}
	roll := rand.Intn(20)
	if roll >= 8 {
		fmt.Println("You successfully escaped!")
		myConfig.EncounteredPokemon = nil
		return nil
	}
	fmt.Println("You were unable to escape!")
	// add attack from wild pokemon here
	return nil
}

func commandAttack(myConfig *gamedata.Config, _ []string) error {
	if myConfig.EncounteredPokemon == nil {
		fmt.Println("You are not in an encounter.")
	}
	if myConfig.PartyPokemon[1] == nil {
		fmt.Println("You have no pokemon in Slot 1!")
		return nil
	}
	mine := myConfig.PartyPokemon[1]
	theirs := myConfig.EncounteredPokemon
	PokemonAttack(mine, theirs)
	PokemonAttack(theirs, mine)
	if theirs.Fainted {
		myConfig.EncounteredPokemon = nil
	}
	return nil
}
