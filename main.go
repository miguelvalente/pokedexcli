package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func()
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays available commands",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandHelp() {
	for _, cmd := range getCommands() {
		fmt.Printf("\t%s: %s\n", cmd.name, cmd.description)
	}
}

func commandExit() {
	fmt.Println("Exiting the Pokedex...")
	os.Exit(0)
}

// func main() {
// 	commands := getCommands()
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for {
// 		fmt.Print("Pokedex > ")
// 		if scanner.Scan() {
// 			input := scanner.Text()
// 			if cmd, found := commands[input]; found {
// 				cmd.callback()
// 			} else {
// 				fmt.Println("Unknown command:", input)
// 			}
// 		}
// 	}
// }

func main() {
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Pokedex!")
	fmt.Println("Usage:")
	commands["help"].callback()
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			if cmd, found := commands[input]; found {
				cmd.callback()
			} else {
				fmt.Println("Unknown command: ", input)
			}
		}
	}
}
