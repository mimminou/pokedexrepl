package main

import (
	"errors"
	"fmt"
	"github.com/mimminou/pokedexrepl/internal/networking"
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
			usage:       "type 'help'",
			function:    printHelp,
		},
		"exit": {
			name:        "exit",
			description: "exits the Pokedex",
			usage:       "type 'exit'",
			function:    exit,
		},
		"map": {
			name:        "map",
			description: "shows next 20 regions in map",
			usage:       "type 'map'",
			function:    mapCmd,
		},
		"mapb": {
			name:        "mapb",
			description: "shows previous 20 regions in map",
			usage:       "type 'mapb'",
			function:    mapBCmd,
		},
		"explore": {
			name:        "explore",
			description: "returns a list of pokemons in a given area",
			usage:       "type 'explore area_name', replace area_name with any valid name returned from the map or mapb, requires map to be called at least once before explore can be used",
			function:    exploreCmd,
		},
		"catch": {
			name:        "catch",
			description: "tries to catch a pokemon",
			usage:       "type 'catch pokemon_name', the pokemon has to be located in the current area of the player",
			function:    catchCmd,
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
	if len(input) < 2 { //check only if 1 arg was passed, nothing more should pass
		val.function()
		return
	}
	val.function(input[1])
}

// define command funcions
func exit(...string) error {
	os.Exit(0)
	return nil
}

func printHelp(...string) error {
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("name : %s  |  description : %s |  usage : %s \n", cmd.name, cmd.description, cmd.usage)
	}
	fmt.Println("")
	return nil
}

var locationArea networking.LocationArea         // this variable holds the current instance of Location area, needed in order to "track" where we are in the pagination
var caughtPokemons map[string]networking.Pokemon //map that holds caught pokemons, indexed by their name

func mapCmd(...string) error {
	endpoint := "/location-area"
	if locationArea.Next == nil && locationArea.Count != 0 {
		fmt.Println("Reached last area, next areas not avaialable")
		return errors.New("Next Area not available")
	}
	if locationArea.Next != nil {
		formattedURL, found := strings.CutPrefix(*locationArea.Next, networking.BaseURL)
		if !found {
			return errors.New("BaseURL not in URL of API call, did the URL change ?")
		} else {
			endpoint = formattedURL
		}
	}

	location, err := networking.GetLocationAreas(endpoint)
	if err != nil {
		fmt.Println(err)
		return err
	}
	locationArea = location

	for _, loc := range locationArea.Results {
		fmt.Println("- " + loc.Name)
	}
	return nil
}

func mapBCmd(...string) error {
	endpoint := "/location-area"
	if locationArea.Count == 0 {
		fmt.Println("No areas loaded, please use the 'map' command first")
		return errors.New("No areas loaded")
	}
	if locationArea.Previous == nil {
		fmt.Println("Reached First Area, no Previous areas available")
		return errors.New("Previous area not available")
	}

	formattedURL, found := strings.CutPrefix(*locationArea.Previous, networking.BaseURL)
	if !found {
		return errors.New("BaseURL not in URL of API call, did the URL change ?")
	} else {
		endpoint = formattedURL
	}

	location, err := networking.GetLocationAreas(endpoint)
	if err != nil {
		fmt.Println(err)
		return err
	}
	locationArea = location
	for _, loc := range locationArea.Results {
		fmt.Println("- " + loc.Name)
	}
	return nil
}

func exploreCmd(area ...string) error {
	if len(area) == 0 {
		fmt.Printf("please pass an area name in this format : 'explore area_name' \n")
		return errors.New("No area specified")
	}
	endpoint := fmt.Sprintf("/location-area/%s", area[0])
	if locationArea.Count == 0 {
		fmt.Println("Please use the 'map' command first to load the areas")
		return errors.New("Areas Not Loaded")
	}

	areaDetails, err := networking.ExploreArea(endpoint)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Exploring . . .")

	for _, pokemon := range areaDetails.PokemonEncounters {
		fmt.Println("- " + pokemon.Pokemon.Name)
	}
	return nil
}

func catchCmd(pokemon ...string) error {
	return nil
}

// define commands
type command struct {
	name        string
	description string
	usage       string
	function    func(...string) error //use variadic because I don't know if the next assignment part is gonig to require further func arg alteration (spoiler, it doesn't)
}
