package main

import (
	"bufio"
	"fmt"
	"iconians/pokedexcli/api"
	"iconians/pokedexcli/internals/pokecache"
	"os"
	"strings"
	"time"
)

type Config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandRegistry map[string]cliCommand
var cfg = &api.Config{
	Cache: pokecache.NewCache(5 * time.Second),
}

func initCommands() {
	commandRegistry = map[string]cliCommand{
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
			description: "Explore the Pokemon world (next 20 locations)",
			callback: func() error {
				api.MapCommand(cfg)
				return nil

			},
		},
		"mapb": {
			name:        "mapb",
			description: "Go back in the Pokemon world (previous 20 locations)",
			callback: func() error {
				api.MapBackCommand(cfg)
				return nil
			},
		},
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range commandRegistry {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func cleanInput(text string) []string {
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)
	if text == "" {
		return []string{}
	}
	return strings.Fields(text)
	// return []string{}
}

func main() {
	initCommands()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Pokedex > ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		// input := strings.TrimSpace(scanner.Text())
		// if input == "" {
		// 	continue
		// }

		words := strings.Fields(scanner.Text())
		if len(words) == 0 {
			continue
			// fmt.Printf("Your command was: %s\n", strings.ToLower(words[0]))
		}

		commandName := words[0]

		cmd, ok := commandRegistry[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
