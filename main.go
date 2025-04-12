package main

import (
	"bufio"
	"fmt"
	"iconians/pokedexcli/api"
	"iconians/pokedexcli/internals/pokecache"
	"math/rand"
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
	callback    func(args ...string) error
}

var commandRegistry map[string]cliCommand
var caughtPokemon = make(map[string]api.Pokemon)

var cfg = &api.Config{
	Cache: pokecache.NewCache(5 * time.Second),
}

func initCommands() {
	commandRegistry = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback: func(args ...string) error {
				return commandExit()
			},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func(args ...string) error {
				return commandHelp()
			},
		},
		"map": {
			name:        "map",
			description: "Explore the Pokemon world (next 20 locations)",
			callback: func(args ...string) error {
				api.MapCommand(cfg)
				return nil

			},
		},
		"mapb": {
			name:        "mapb",
			description: "Go back in the Pokemon world (previous 20 locations)",
			callback: func(args ...string) error {
				api.MapBackCommand(cfg)
				return nil
			},
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area. Usage: explore <area-name>",
			callback: func(args ...string) error {
				if len(args) < 1 {
					return fmt.Errorf("please provide a location area name")
				}
				return api.ExploreArea(cfg, args[0])
			},
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon. Usage: catch <pokemon-name>",
			callback: func(args ...string) error {
				return catch(args...)
			},
		},
		"inspect": {
			name:        "inpect",
			description: "Inspect a caught Pokemon. Usage: inspect <pokemon-name>",
			callback: func(args ...string) error {
				return inspectCommand(args...)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught Pokemon",
			callback: func(args ...string) error {
				return pokedexCommand()
			},
		},
	}
}

func catch(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a pokemon name")
	}

	name := args[0]

	if _, caught := caughtPokemon[name]; caught {
		fmt.Printf("%s is already in your Pokedex!\n", name)
		return nil
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := api.GetPokemon(cfg, name)
	if err != nil {
		return err
	}

	chance := 100 - pokemon.BaseExperience
	if chance < 10 {
		chance = 10
	}

	if rand.Intn(100) < chance {
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
		caughtPokemon[name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func inspectCommand(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a Pokemon name. Usage: inspect <pokemon-name>")
	}

	name := args[0]
	pokemon, ok := caughtPokemon[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.StatInfo.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func pokedexCommand(args ...string) error {
	if len(caughtPokemon) == 0 {
		fmt.Println("You haven't caught any Pokemon yet.")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range caughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
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
		args := words[1:]

		cmd, ok := commandRegistry[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(args...); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
