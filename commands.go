package main

import (
	"fmt"
	"os"
	"strings"
)

var commands map[string]command

func init() {
	// map commands to input str
	commands = map[string]command{
		"help": {
			name:        "help",
			description: "prints this help message",
			function:    printHelp,
		},
		"exit": {
			name:        "exit",
			description: "exits the Pokedex",
			function:    exit,
		},
		"map": {
			name:        "map",
			description: "show next 20 regions in map",
			function:    mapCmd,
		},
		"mapb": {
			name:        "mapb",
			description: "show previous 20 regions in map",
			function:    mapBCmd,
		},
	}
}

func sanitize(input string) []string {
	inputLowered := strings.ToLower(input)
	words := strings.Fields(inputLowered)
	return words
}

func runCommand(input []string) {
	val, ok := commands[input[0]]
	if !ok {
		fmt.Println("command not found")
		return
	}
	val.function()
}

// define command funcions
func exit() error {
	os.Exit(0)
	return nil
}

func printHelp() error {
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("name : %s  |  description : %s \n", cmd.name, cmd.description)
	}
	return nil
}

func mapCmd() error {
	// fill this
	return nil
}

func mapBCmd() error {
	// fill this
	return nil
}

// define commands
type command struct {
	name        string
	description string
	function    func() error
}

var helpCmd command = command{
	name:        "help",
	description: "prints this help message",
	function:    printHelp,
}

var exitCmd command = command{
	name:        "exit",
	description: "exits the pokedex",
	function:    exit,
}
