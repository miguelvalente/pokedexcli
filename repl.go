package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/miguelvalente/pokedexcli/commands"
	"github.com/miguelvalente/pokedexcli/internal/pokecache"
)

func startRepl() {
	commands_ := commands.GetCommands()
	scanner := bufio.NewScanner(os.Stdin)
	config := &commands.Config{}
	const baseTime = 100 * time.Second
	cache := pokecache.NewCache(baseTime)
	pokedex := &commands.Pokedex{
		MyPokemons: make(map[string]*commands.Pokemon),
	}

	fmt.Println("Welcome to Pokedex!")
	fmt.Println("Usage:")
	commands_["help"].Callback("", config, cache, pokedex)
	for {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			input := scanner.Text()

			fields := strings.Fields(input)
			if len(fields) == 0 {
				fmt.Println("Empty input")
				continue
			}
			commandName := fields[0]
			extra := ""
			if len(fields) > 1 {
				extra = fields[1]
			}
			if cmd, found := commands_[commandName]; found {
				cmd.Callback(extra, config, cache, pokedex)
			} else {
				fmt.Println("Unknown command: ", input)
			}
		}
	}
}
