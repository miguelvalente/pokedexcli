package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type AutoGenerated struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap() {
	resp, err := http.Get("https://pokeapi.co/api/v2/location-area/")
	if err != nil {
		fmt.Println("boo")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	data := AutoGenerated{}
	err_unmarsh := json.Unmarshal(body, &data)
	if err_unmarsh != nil {
		fmt.Println("2", err_unmarsh)
	}
	for _, res := range data.Results {
		fmt.Println(res.Name)
	}
}

func commandMapb() {
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
		"map": {
			name:        "map",
			description: "List available areas",
			Callback:    commandMap,
		},
	}
}
