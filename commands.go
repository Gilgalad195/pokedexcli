package main

import "github.com/gilgalad195/pokedexcli/internal/pokeapi"

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.Config) error
}

var commands map[string]cliCommand

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
			description: "Displays next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Diaplays previous 20 location areas",
			callback:    commandMapb,
		},
	}
}
