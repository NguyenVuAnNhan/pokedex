package main

import(
	"fmt";
	"bufio";
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func() error
}

func initCommands() map[string]cliCommand {
	commands := map[string]cliCommand {}
	commands["exit"] = cliCommand{
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	}
	commands["help"] = cliCommand{
		name: "help",
		description: "Displays a help message",
		callback: func() error {
			return commandHelp(commands)
		},
	}
	return commands
}

func main(){
	cliCommands := initCommands()
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

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cliCommands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n")
	for _, cmd := range cliCommands {
		fmt.Printf("  %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}