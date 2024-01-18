package main

import (
	"errors"
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
		"loc": {
			name:        "loc",
			description: "checks location area var",
			function:    checkLocationArea,
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

var locationArea LocationArea // this variable holds the current instance of Location area, needed in order to "track" where we are in the pagination

func mapCmd() error {
	endpoint := "/location-area"
	if locationArea.Next == nil && locationArea.Count != 0 {
		fmt.Println("Reached last area, next areas not avaialable")
		return errors.New("Next Area not available")
	}
	if locationArea.Next != nil {
		fmt.Println("Next URL exists")
		formattedURL, found := strings.CutPrefix(*locationArea.Next, baseURL)
		if !found {
			return errors.New("BaseURL not in URL of API call, did the URL change ?")
		} else {
			endpoint = formattedURL
		}
	}

	location, err := getLocationAreas(endpoint)
	if err != nil {
		fmt.Println(err)
		return err
	}
	locationArea = location
	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func mapBCmd() error {
	endpoint := "/location-area"
	if locationArea.Count == 0 {
		fmt.Println("No areas loaded, please use the 'map' command first")
		return errors.New("No areas loaded")
	}
	if locationArea.Previous == nil {
		fmt.Println("Reached First Area, no Previous areas available")
		return errors.New("Previous area not available")
	}

	formattedURL, found := strings.CutPrefix(*locationArea.Previous, baseURL)
	if !found {
		return errors.New("BaseURL not in URL of API call, did the URL change ?")
	} else {
		endpoint = formattedURL
	}

	location, err := getLocationAreas(endpoint)
	if err != nil {
		fmt.Println(err)
		return err
	}
	locationArea = location
	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func checkLocationArea() error {
	fmt.Println(locationArea.Count)
	return nil
}

// define commands
type command struct {
	name        string
	description string
	function    func() error
}
