package main

import (
	"bufio"
	"fmt"
	"os"
)

func fetchData(pokemon string) (string, error) {
	mockdata := "just some random stuff, this should be structured JSON data then formatted"
	return mockdata, nil
}

func main() {
	fmt.Println("Welcome to the GoGo second generation PokeDex")
	fmt.Println("For help, please type 'help'")
	reader := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()
		input := reader.Text()
		if len(input) == 0 {
			fmt.Println("Please enter a valid input")
			continue
		}
		sanitizedInput := sanitize(input)
		runCommand(sanitizedInput)

	}
}
