package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/commands"
)

func startRepl() {
	commands_ := commands.GetCommands()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to Pokedex!")
	fmt.Println("Usage:")
	commands_["help"].Callback()
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()
			if cmd, found := commands_[input]; found {
				cmd.Callback()
			} else {
				fmt.Println("Unknown command: ", input)
			}
		}
	}
}
