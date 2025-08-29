package main

import (
	"github.com/gilgalad195/pokedexcli/internal/gamedata"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*gamedata.Config, []string) error
	//this is so the commands can receive and update the previous/next state.
}

var commands map[string]cliCommand

// definiing and using init() is to prevent circular dependencies between here and main.
func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays a list of visited locations",
			callback:    commandMap,
		},
		"explore": {
			name:        "explore",
			description: "Searches an area for Pokemon",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a Pokemon that you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all of the pokemon you have caught",
			callback:    commandPokedex,
		},
		"save": {
			name:        "save",
			description: "Saves your game",
			callback:    commandSave,
		},
		"load": {
			name:        "load",
			description: "Loads a previous gamestate",
			callback:    commandLoad,
		},
		"look": {
			name:        "look",
			description: "Looks around to find possible paths",
			callback:    commandLook,
		},
		"move": {
			name:        "move",
			description: "Moves to the desired location",
			callback:    commandMove,
		},
		"party": {
			name:        "party",
			description: "add, remove, or inspect pokemon in party",
			callback:    commandParty,
		},
		"run": {
			name:        "run",
			description: "escape from the current encounter",
			callback:    commandRun,
		},
		"attack": {
			name:        "run",
			description: "attacks the encountered pokemon",
			callback:    commandAttack,
		},
	}
}
