package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/miguelvalente/pokedexcli/internal/pokecache"
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

func commandMap(config *mapConfig, cache *pokecache.Cache) {
	url := "https://pokeapi.co/api/v2/location-area/"

	if config.Next != "" {
		url = config.Next
	}

	body := []byte{}

	value, exists := cache.Get(url)
	if exists {
		body = value
	} else {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		cache.Add(url, body)
	}

	var data MapResponse
	err := json.Unmarshal(body, &data)
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

func commandMapb(config *mapConfig, cache *pokecache.Cache) {
	if config.Previous == "" {
		fmt.Println("No previous location to visit")
		return
	}

	body, found := cache.Get(config.Previous)
	if !found {
		// fmt.Println("Cache miss for key:", config.Previous)

		resp, err := http.Get(config.Previous)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}
		cache.Add(config.Previous, body)
	}

	var data MapResponse
	if err := json.Unmarshal(body, &data); err != nil {
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

func commandMapFunc(config *mapConfig, cache *pokecache.Cache) func() {
	return func() {
		commandMap(config, cache)
	}
}

func commandMapbFunc(config *mapConfig, cache *pokecache.Cache) func() {
	return func() {
		commandMapb(config, cache)
	}
}
func GetCommands() map[string]cliCommand {
	config := &mapConfig{}
	const baseTime = 100 * time.Second
	cache := pokecache.NewCache(baseTime)

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
			Callback:    commandMapFunc(config, cache),
		},
		"mapb": {
			name:        "mapb",
			description: "List previous available areas",
			Callback:    commandMapbFunc(config, cache),
		},
	}
}
