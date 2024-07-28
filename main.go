package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Config struct {
	Next     *string
	Previous *string
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
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

func commandMap(config *Config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.Next != nil {
		url = *config.Next
	}
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data: %s", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	var locationResp LocationAreaResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return err
	}
	fmt.Println("Location areas:")
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}

	config.Next = locationResp.Next
	config.Previous = locationResp.Previous
	return nil
}

func commandMapb(config *Config) error {
	if config.Previous == nil {
		return fmt.Errorf("you are already at the first page")
	}
	url := *config.Previous
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("response failed with status code: %d and body: %s", res.StatusCode, body)
	}
	var locationResp LocationAreaResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return err
	}
	fmt.Println("Location areas:")
	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}

	config.Next = locationResp.Next
	config.Previous = locationResp.Previous
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{}
	cmd := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		cmd, ok := cmd[input]
		if ok {
			err := cmd.callback(config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
			}
		} else {
			fmt.Println("Command not found. Type 'help' to see the list of available commands.")
		}
	}
}
