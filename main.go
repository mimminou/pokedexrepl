package main

import (
	"bufio"
	"fmt"
	"os"
)

// define commands
var commands []command

var helpCmd command = command{
	name:        "/help",
	description: "prints this help message",
}
var exit command = command{
	name:        "/exit",
	description: "exits the pokedex",
}

type command struct {
	name        string
	description string
}

func interpretCommand(cmd string) {
	//switch that tests for commands
	switch cmd {
	case "/help":
		fmt.Println("print help here")
		printHelp()

	case "/exit":
		fmt.Println("Exiting. . .")
		os.Exit(0)

	default:
		returnedData, fetchErr := fetchData(cmd)
		if fetchErr != nil {
			fmt.Println(returnedData)
		} else {
			fmt.Println(fetchErr)
		}
	}
}

func printHelp() error {
	return nil
}

func fetchData(pokemon string) (string, error) {
	mockdata := "just some random stuff, this should be structured JSON data then formatted"
	return mockdata, nil
}

// todo : create fetching function, create help command and structs
func main() {

	fmt.Println("Welcome to the GoGo second generation PokeDex")
	fmt.Println("For help, please type /help")
	reader := bufio.NewScanner(os.Stdin)

	for {
		for reader.Scan() {
			input := reader.Text()
			interpretCommand(input)
		}
	}

}
