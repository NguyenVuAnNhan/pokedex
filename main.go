package main

import(
	"fmt"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
	pokecache "github.com/NguyenVuAnNhan/pokedexcli/internal/pokecache"
	"time"
	"io"
	"math/rand"
)

type cliCommand struct {
	name string
	description string
	callback func(args []string) error
}

type config struct {
	Next *string
	Previous *string
	Cache *pokecache.Cache
	Pokedex Pokedex
}

type apiResponse[T any] struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
	Results []T `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type PokemonList struct {
    PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
    Pokemon        PokemonEntry `json:"pokemon"`
    VersionDetails []any   `json:"version_details"` // ignore for now
}

type PokemonEntry struct {
	Name string `json:"name"`
	URL string `json:"url"`
}

type Stat struct {
	BaseStat int `json:"base_stat"`
	StatEntry StatEntry `json:"stat"`
}

type StatEntry struct {
	StatName string `json:"name"`
}

type Type struct {
	TypeEntry TypeEntry `json:"type"`
}

type TypeEntry struct {
	TypeName string `json:"name"`
}

type Pokemon struct {
	Name string `json:"name"`
	Height int `json:"height"`
	Weight int `json:"weight"`
	Stats []Stat `json:"stats"`
	Types []Type `json:"types"`
	BaseExperience int `json:"base_experience"`
}

type Pokedex struct {
	Caught map[string]Pokemon
}

func initCommands(cfg *config) map[string]cliCommand {
	commands := map[string]cliCommand {}
	commands["exit"] = cliCommand{
		name: "exit",
		description: "Exit the Pokedex",
		callback: func(args []string) error {
			return commandExit(cfg)
		},
	}
	commands["help"] = cliCommand{
		name: "help",
		description: "Displays a help message",
		callback: func(args []string) error {
			return commandHelp(commands, cfg)
		},
	}
	commands["map"] = cliCommand{
		name: "map",
		description: "Displays the next page of the map of the Pokemon world",
		callback: func(args []string) error {
			return commandMap(cfg)
		},
	}
	commands["mapb"] = cliCommand{
		name: "mapb",
		description: "Displays the previous page of the map of the Pokemon world",
		callback: func(args []string) error {
			return commandMapBack(cfg)
		},
	}
	commands["explore"] = cliCommand{
		name: "explore",
		description: "Displays the next pokemons in the region",
		callback: func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing target location area")
			}
			return commandExplore(cfg, args[0])
		},
	}
	commands["catch"] = cliCommand{
		name: "catch",
		description: "Catch a pokemon and add it to your pokedex",
		callback: func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing target pokemon")
			}
			return commandCatch(cfg, args[0])
		},
	}
	commands["inspect"] = cliCommand{
		name: "inspect",
		description: "Inspect a pokemon in your pokedex",
		callback: func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missing target pokemon")
			}
			pokemon, exists := cfg.Pokedex.Caught[args[0]]
			if !exists {
				fmt.Printf("You haven't caught %s yet!\n", args[0])
				return nil
			}
			fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nBase Experience: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight, pokemon.BaseExperience)
			fmt.Println("Stats:")
			for _, stat := range pokemon.Stats {
				fmt.Printf(" - %s: %d\n", stat.StatEntry.StatName, stat.BaseStat)
			}
			fmt.Println("Types:")
			for _, t := range pokemon.Types {
				fmt.Printf(" - %s\n", t.TypeEntry.TypeName)
			}
			return nil
		},
	}
	return commands
}

func main(){
	startURL := "https://pokeapi.co/api/v2/location-area"

	cache := pokecache.NewCache(5 * time.Second)

	pokedex := Pokedex{
		Caught: make(map[string]Pokemon),
	}

	cfg := &config{
		Next: &startURL,
		Previous: nil,
		Cache: cache,
		Pokedex: pokedex,
	}
	cliCommands := initCommands(cfg)

	var scanner = bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		command := cleanInput(input)[0]
		args := cleanInput(input)[1:]

		cliCommand, exists := cliCommands[command]
		if exists {
			err := cliCommand.callback(args)
			if err != nil {
				fmt.Printf("Error executing command: %v\n", err)
			}
		} else {
		fmt.Println("Unknown command")
		}
	}
}

func commandExit(_cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cliCommands map[string]cliCommand, _cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, cmd := range cliCommands {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	if cfg.Next == nil {
		fmt.Println("No more locations to display")
		return nil
	}

	var resp apiResponse[locationArea]
	var res *http.Response
	var err error
	var data []byte
	url := *cfg.Next

	if cached, exists := cfg.Cache.Get(url); exists {
		// fmt.Println("cache hit:", url)
		data = cached
	} else {
		// fmt.Println("cache miss:", url)
		res, err = http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %s", res.Status)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, data)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	cfg.Previous = resp.Previous
	cfg.Next = resp.Next

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	var resp apiResponse[locationArea]
	var res *http.Response
	var data []byte
	var err error
	url := *cfg.Previous

	if cached, exists := cfg.Cache.Get(url); exists {
		// fmt.Println("cache hit:", url)
		data = cached
	} else {
		// fmt.Println("cache miss:", url)
		res, err = http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %s", res.Status)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, data)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return err
	}

	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	cfg.Previous = resp.Previous
	cfg.Next = resp.Next

	return nil
}

func commandExplore(cfg *config, target string) error {
	var res *http.Response
	var data []byte
	var err error
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", target)

	fmt.Println("Exploring", target + "...")

	if cached, exists := cfg.Cache.Get(url); exists {
		// fmt.Println("cache hit:", url)
		data = cached
	} else {
		// fmt.Println("cache miss:", url)
		res, err = http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %s", res.Status)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, data)
	}

	var pokemonList PokemonList
	err = json.Unmarshal(data, &pokemonList)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, pokemonEncounter := range pokemonList.PokemonEncounters {
		fmt.Println(" -", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, target string) error {
	var res *http.Response
	var data []byte
	var err error
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", target)

	if cached, exists := cfg.Cache.Get(url); exists {
		// fmt.Println("cache hit:", url)
		data = cached
	} else {
		// fmt.Println("cache miss:", url)
		res, err = http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status: %s", res.Status)
		}

		data, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cfg.Cache.Add(url, data)
	}

	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if rand.Float64()*700 > float64(pokemon.BaseExperience) {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.Pokedex.Caught[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}