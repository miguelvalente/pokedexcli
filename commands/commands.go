package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/miguelvalente/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	Callback    func(string, *Config, *pokecache.Cache, *Pokedex)
}

type Config struct {
	Next     string
	Previous string
}

func CommandHelp(input string, cfg *Config, cache *pokecache.Cache, pk *Pokedex) {
	for _, cmd := range GetCommands() {
		fmt.Printf("\t%s: %s\n", cmd.name, cmd.description)
	}
}

func commandExit(input string, cfg *Config, cache *pokecache.Cache, pk *Pokedex) {
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

func commandMap(input string, config *Config, cache *pokecache.Cache, pk *Pokedex) {
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

func commandMapb(input string, config *Config, cache *pokecache.Cache, pk *Pokedex) {
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

func commandExplore(input string, config *Config, cache *pokecache.Cache, pk *Pokedex) {
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

type Pokedex struct {
	MyPokemons map[string]*Pokemon
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func commandCatch(input string, config *Config, cache *pokecache.Cache, pokedex *Pokedex) {

	url := "https://pokeapi.co/api/v2/pokemon/" + input

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

	var data Pokemon
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("Throwing a Pokeball at ", input, "... ")

	if rand.Intn(data.BaseExperience) > data.BaseExperience-40 {
		fmt.Println(input, " was caught!")
		pokedex.MyPokemons[input] = &data
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Println(input, " escaped!")
	}
}

func commandInspect(input string, config *Config, cache *pokecache.Cache, pokedex *Pokedex) {
	if pokemon, ok := pokedex.MyPokemons[input]; ok {
		fmt.Println("Name: ", input)
		fmt.Println("Height: ", pokemon.Height)
		fmt.Println("Weight: ", pokemon.Weight)
		fmt.Println("Stats: ")
		for _, stat := range pokemon.Stats {
			fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types: ")

		for _, pokeType := range pokemon.Types {
			fmt.Println(" -", pokeType.Type.Name)

		}

	} else {
		fmt.Println("you have not caught that pokemon")
	}
}

func commandPokedex(input string, config *Config, cache *pokecache.Cache, pokedex *Pokedex) {
	fmt.Println("Your Pokedex:")
	for pokemonName := range pokedex.MyPokemons {
		fmt.Println(" -", pokemonName)
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
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			Callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caugth pokemon",
			Callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Print pokedex",
			Callback:    commandPokedex,
		},
	}
}
