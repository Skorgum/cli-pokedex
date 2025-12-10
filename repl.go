package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Skorgum/cli-pokedex/internal/pokeapi"
)

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}
func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		PokeClient: pokeapi.NewClient(5 * time.Second),
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}
		cmdName := words[0]
		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}
		cmd, ok := commands[cmdName]
		if !ok {
			fmt.Printf("Unknown command: %s\n", cmdName)
			continue
		}
		err := cmd.callback(cfg, args...)
		if err != nil {
			fmt.Printf("Error executing command %s: %s\n", cmdName, err)
			continue
		}
	}
}
func commandExit(cfg *config, args ...string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	order := []string{"help", "exit", "map", "mapb", "explore"}

	for _, name := range order {
		cmd := commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	res, err := cfg.PokeClient.ListLocations(cfg.Next)
	if err != nil {
		return err
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous

	return nil
}

func commandMapB(cfg *config, args ...string) error {
	if cfg.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	res, err := cfg.PokeClient.ListLocations(cfg.Previous)
	if err != nil {
		return err
	}

	for _, location := range res.Results {
		fmt.Println(location.Name)
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous

	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("Please specify an area")
	}

	locationName := args[0]
	fmt.Printf("Exploring %s...\n", locationName)

	loc, err := cfg.PokeClient.GetLocation(locationName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokémon:")
	for _, encounter := range loc.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

type config struct {
	Next       *string
	Previous   *string
	PokeClient pokeapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
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
			description: "Displays 20 locations. Subsequent call displays the next 20, and so on",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area to find Pokémon",
			callback:    commandExplore,
		},
	}
}
