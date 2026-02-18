package main

import(
	"fmt";
	"bufio";
	"os";
	"net/http";
	"encoding/json";
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

type config struct {
	Next *string
	Previous *string
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

func initCommands(cfg *config) map[string]cliCommand {
	commands := map[string]cliCommand {}
	commands["exit"] = cliCommand{
		name: "exit",
		description: "Exit the Pokedex",
		callback: func() error {
			return commandExit(cfg)
		},
	}
	commands["help"] = cliCommand{
		name: "help",
		description: "Displays a help message",
		callback: func() error {
			return commandHelp(commands, cfg)
		},
	}
	commands["map"] = cliCommand{
		name: "map",
		description: "Displays the next page of the map of the Pokemon world",
		callback: func() error {
			return commandMap(cfg)
		},
	}
	commands["mapb"] = cliCommand{
		name: "mapb",
		description: "Displays the previous page of the map of the Pokemon world",
		callback: func() error {
			return commandMapBack(cfg)
		},
	}
	return commands
}

func main(){
	startURL := "https://pokeapi.co/api/v2/location-area"

	cfg := &config{
		Next: &startURL,
		Previous: nil,
	}
	cliCommands := initCommands(cfg)

	var scanner = bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		command := cleanInput(input)[0]
		cliCommand, exists := cliCommands[command]
		if exists {
			err := cliCommand.callback()
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
	fmt.Println("Usage:\n")
	for _, cmd := range cliCommands {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	var resp apiResponse[locationArea]

	if cfg.Next == nil {
		fmt.Println("No more locations to display")
		return nil
	}

	res, err := http.Get(*cfg.Next)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}

	// fmt.Println(resp)

	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	if resp.Previous != nil {
		cfg.Previous = resp.Previous
	} else {
		cfg.Previous = nil
	}
	if resp.Next != nil {
		cfg.Next = resp.Next
	} else {
		cfg.Next = nil
	}
	return nil
}

func commandMapBack(cfg *config) error {
	var resp apiResponse[locationArea]

	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(*cfg.Previous)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return err
	}

	// fmt.Println(resp)

	for _, location := range resp.Results {
		fmt.Println(location.Name)
	}

	if resp.Previous != nil {
		cfg.Previous = resp.Previous
	} else {
		cfg.Previous = nil
	}
	if resp.Next != nil {
		cfg.Next = resp.Next
	} else {
		cfg.Next = nil
	}
	return nil
}