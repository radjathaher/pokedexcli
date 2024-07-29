package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"pokedexcli/internal/pokeapi"
	"strings"
	"time"
)

type Config struct {
	Next     *string
	Previous *string
	Client   *pokeapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Prints the list of available commands",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display names of the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display names of the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore the location areas",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
	}
}

func commandHelp(cfg *Config, args ...string) error {
	commands := getCommands()
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *Config, args ...string) error {
	fmt.Println("Exiting program...")
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config, args ...string) error {
	res, err := cfg.Client.GetLocationAreas(cfg.Next)
	if err != nil {
		return err
	}
	fmt.Println("Location areas:")
	for _, area := range res.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous
	return nil
}

func commandMapb(cfg *Config, args ...string) error {
	if cfg.Previous == nil {
		return fmt.Errorf("you are already at the first page")
	}

	res, err := cfg.Client.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	fmt.Println("Location areas:")
	for _, area := range res.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = res.Next
	cfg.Previous = res.Previous
	return nil
}

func commandExplore(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a location area name to explore")
	}
	locationAreaName := args[0]
	locationArea, err := cfg.Client.GetLocationArea(locationAreaName)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", locationAreaName)
	fmt.Println("Found Pokemon:")
	for _, encounter := range locationArea.PokemonEncounters {
		fmt.Printf("%s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a pokemon name to catch")
	}
	pokemonName := args[0]
	fmt.Printf("Throwing a pokeball at %s...\n", pokemonName)
	pokemonData, err := cfg.Client.GetPokemonData(pokemonName)
	if err != nil {
		return fmt.Errorf("error fetching pokemon data")
	}
	if rand.Intn(pokemonData.BaseExperience)*2 > pokemonData.BaseExperience {
		fmt.Printf("%s was caught\n", pokemonName)
		cfg.Client.AddToPokedex(pokemonData)
		return nil
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}
}

func commandInspect(cfg *Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a pokemon name to inspect")
	}
	pokemonName := args[0]
	pokemon, ok := cfg.Client.GetFromPokedex(pokemonName)
	if !ok {
		return fmt.Errorf("you have not caught that pokemon yet")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &Config{
		Client: pokeapi.NewClient(10 * time.Second),
	}
	cmd := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		words := strings.Fields(input)
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		args := words[1:]

		cmd, ok := cmd[commandName]
		if ok {
			err := cmd.callback(cfg, args...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
			}
		} else {
			fmt.Println("Command not found. Type 'help' to see the list of available commands.")
		}
	}
}
