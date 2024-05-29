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

type mapConfig struct {
	Next     string
	Previous string
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

type MapResponse struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(config *mapConfig) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if config.Next != "" {
		url = config.Next
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var data MapResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, res := range data.Results {
		fmt.Println(res.Name)
	}

	config.Next = data.Next
	if data.Previous != nil {
		config.Previous = *data.Previous
	} else {
		config.Previous = ""
	}
}

func commandMapb(config *mapConfig) {
	if config.Previous == "" {
		fmt.Println("No previous location to visit")
		return
	}

	resp, err := http.Get(config.Previous)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var data MapResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	for _, res := range data.Results {
		fmt.Println(res.Name)
	}

	// config.Previous = *data.Previous
	// config.Next = data.Next
	// Update the Config with the Next and Previous URLs
	config.Next = data.Next
	if data.Previous != nil {
		config.Previous = *data.Previous
	} else {
		config.Previous = ""
	}

}

func commandMapFunc(config *mapConfig) func() {
	return func() {
		commandMap(config)
	}
}

func commandMapbFunc(config *mapConfig) func() {
	return func() {
		commandMapb(config)
	}
}
func GetCommands() map[string]cliCommand {
	config := &mapConfig{}

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
			description: "List nexy available areas",
			Callback:    commandMapFunc(config),
		},
		"mapb": {
			name:        "mapb",
			description: "List previous available areas",
			Callback:    commandMapbFunc(config),
		},
	}
}
