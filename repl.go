package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Skorgum/cli-pokedex/internal/pokeapi"
)

type config struct {
	Next       *string
	Previous   *string
	PokeClient pokeapi.Client
	Pokedex    map[string]pokeapi.Pokemon
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch indicated Pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect indicated Pokémon",
			callback:    commandInspect,
		},
	}
}

func cleanInput(text string) []string {
	lower := strings.ToLower(text)
	words := strings.Fields(lower)
	return words
}
func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &config{
		PokeClient: pokeapi.NewClient(5 * time.Second),
		Pokedex:    make(map[string]pokeapi.Pokemon),
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

	order := []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect"}

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
		return fmt.Errorf("please specify an area")
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a Pokémon")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemonData, err := cfg.PokeClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	baseExp := pokemonData.BaseExperience

	var catchChance int
	switch {
	case baseExp >= 300:
		catchChance = 25
	case baseExp >= 100:
		catchChance = 50
	default:
		catchChance = 75
	}

	roll := rand.Intn(100)
	if roll < catchChance {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.Pokedex[pokemonName] = pokemonData
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a Pokémon")
	}

	pokemonName := args[0]

	p, ok := cfg.Pokedex[pokemonName]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)

	fmt.Println("Stats:")
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}

	return nil
}
