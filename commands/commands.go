package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/miguelvalente/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(string, *MapConfig, *pokecache.Cache)
}

type MapConfig struct {
	Next     string
	Previous string
}

func CommandHelp(input string, cfg *MapConfig, cache *pokecache.Cache) {
	for _, cmd := range GetCommands() {
		fmt.Printf("\t%s: %s\n", cmd.name, cmd.description)
	}
}

func commandExit(input string, cfg *MapConfig, cache *pokecache.Cache) {
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

func commandMap(input string, config *MapConfig, cache *pokecache.Cache) {
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

func commandMapb(input string, config *MapConfig, cache *pokecache.Cache) {
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

type PokemonResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandExplore(input string, config *MapConfig, cache *pokecache.Cache) {
	url := "https://pokeapi.co/api/v2/location-area/" + input

	body, found := cache.Get(url)
	if !found {
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

	var data PokemonResponse
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("Exploring ", input, "...")
	if len(data.PokemonEncounters) >= 1 {
		fmt.Println("Found Pokemon:")
		for _, res := range data.PokemonEncounters {
			fmt.Println(res.Pokemon.Name)
		}

	} else {
		fmt.Println("No Pokemons Found")
	}
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
			description: "List next available areas",
			Callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "List previous available areas",
			Callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Lists pokemons in the area",
			Callback:    commandExplore,
		},
	}
}
