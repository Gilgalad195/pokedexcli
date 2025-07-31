package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gilgalad195/pokedexcli/internal/pokeapi"
	"github.com/gilgalad195/pokedexcli/internal/pokecache"
)

func main() {

	pokeCache := pokecache.NewCache(5 * time.Second)

	myConfig := &pokeapi.Config{
		//this initializes the struct with the starting API url.
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Cache:    pokeCache,
	}

	//this is my REPL loop, which looks for user input and executes commands
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		hasToken := scanner.Scan()
		if hasToken {
			lowerString := cleanInput(scanner.Text())
			userInput := lowerString[0]
			if command, exists := commands[userInput]; exists {
				err := command.callback(myConfig)
				if err != nil {
					fmt.Println("Error:", err)
				}
			} else {
				fmt.Println("Unknown command")
			}
		}
	}
}

func cleanInput(text string) []string {
	cleanText := strings.Fields(strings.ToLower(text))
	return cleanText
}

func commandExit(_ *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(_ *pokeapi.Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(myConfig *pokeapi.Config) error {
	if myConfig.Next == "" {
		fmt.Println("you're on the last page")
	} else {
		//declaring body and err for use later
		var body []byte
		var err error
		val, found := myConfig.Cache.Get(myConfig.Next)
		if found {
			fmt.Println("This was in cache!")
			body = val
		} else {
			body, err = pokeapi.GetLocations(myConfig.Next, myConfig)
			if err != nil {
				return fmt.Errorf("failed to get locations: %v", err)
			}
			myConfig.Cache.Add(myConfig.Next, body)
		}

		data, err := pokeapi.FormatResponse(body)
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

func commandMapb(myConfig *pokeapi.Config) error {
	if myConfig.Previous == "" {
		fmt.Println("you're on the first page")
	} else {
		var body []byte
		var err error
		val, found := myConfig.Cache.Get(myConfig.Previous)
		if found {
			fmt.Println("This was in cache!")
			body = val
		} else {
			body, err = pokeapi.GetLocations(myConfig.Previous, myConfig)
			if err != nil {
				return fmt.Errorf("failed to get locations: %v", err)
			}
			myConfig.Cache.Add(myConfig.Previous, body)
		}

		data, err := pokeapi.FormatResponse(body)
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
