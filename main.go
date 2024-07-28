package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
)

type Config struct {
	Next     *string
	Previous *string
	Client   *pokeapi.Client
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Prints the list of available commands",
			callback:    func(*Config) error { return commandHelp() },
		},
		"exit": {
			name:        "exit",
			description: "Exits the program",
			callback:    func(*Config) error { return commandExit() },
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
	}
}

func commandHelp() error {
	commands := getCommands()
	fmt.Println("Available commands:")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit() error {
	fmt.Println("Exiting program...")
	os.Exit(0)
	return nil
}

func commandMap(cfg *Config) error {
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

func commandMapb(cfg *Config) error {
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := &Config{
		Client: pokeapi.NewClient(),
	}
	cmd := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		cmd, ok := cmd[input]
		if ok {
			err := cmd.callback(cfg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
			}
		} else {
			fmt.Println("Command not found. Type 'help' to see the list of available commands.")
		}
	}
}
