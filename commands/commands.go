package commands

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	Callback    func()
}

func CommandHelp() {
	for _, cmd := range GetCommands() {
		fmt.Printf("\t%s: %s\n", cmd.name, cmd.description)
	}
}

func commandExit() {
	fmt.Println("Exiting the Pokedex...")
	os.Exit(0)
}

func GetCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays available commands",
			Callback:    CommandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			Callback:    commandExit,
		},
	}
}
