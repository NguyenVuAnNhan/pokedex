package main

import(
	"fmt";
	"bufio";
	"os"
)

func main(){
	var scanner = bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		clean_input := cleanInput(input)
		fmt.Println("Your command was:", clean_input[0])
	}
}	