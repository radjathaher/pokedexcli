package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cmd := getCommands()
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()

		cmd, ok := cmd[input]
		if ok {
			err := cmd.callback()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %s\n", err)
			}
		} else {
			fmt.Println("Command not found. Type 'help' to see the list of available commands.")
		}
	}
}
