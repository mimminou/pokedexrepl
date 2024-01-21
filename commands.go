package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mimminou/pokedexrepl/internal/networking"
)

var commands map[string]command

func init() {
	// map commands to input str
	commands = map[string]command{
		"help": {
			name:        "help",
			description: "Prints this help message",
			usage:       "Type 'help'",
			function:    helpCmd,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			usage:       "Type 'exit'",
			function:    exitCmd,
		},
		"map": {
			name:        "map",
			description: "Shows the next 20 regions in map",
			usage:       "Type 'map'",
			function:    mapCmd,
		},
		"mapb": {
			name:        "mapb",
			description: "Shows the previous 20 regions in map",
			usage:       "Type 'mapb'",
			function:    mapBCmd,
		},
		"explore": {
			name:        "explore",
			description: "Returns a list of Pokemons in a given area",
			usage:       "Type 'explore area_name', replace area_name with any valid name returned from the map or mapb commands, requires map to be used at least once before explore can be used",
			function:    exploreCmd,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a pokemon, once caught, it will be added to the Pokedex",
			usage:       "Type 'catch pokemon_name', the pokemon has to be located in the current area of the player",
			function:    catchCmd,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspects a Pokemon in the Pokedex",
			usage:       "type 'inspect pokemon_name', the pokemon has to be already caught by the player",
			function:    inspectCmd,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the Pokemons in the Pokedex",
			usage:       "type 'pokedex'",
			function:    pokedexCmd,
		},
	}
	caughtPokemons = make(map[string]networking.Pokemon) // init the map
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
func exitCmd(...string) error {
	os.Exit(0)
	return nil
}

func helpCmd(...string) error {
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("Command : %s\n", cmd.name)
		fmt.Printf("  usage: %s\n", cmd.usage)
		fmt.Printf("  description : %s\n\n", cmd.description)

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
	if len(pokemon) == 0 {
		fmt.Printf("please pass a Pokemon name in this format : 'catch pokemon_name' \n")
		return errors.New("No Pokemon name specified")
	}
	if entry, exists := caughtPokemons[pokemon[0]]; exists {
		pokemonDetails := entry
		isCatched := calculateCatchProbability(int32(pokemonDetails.BaseExperience))
		if isCatched {
			caughtPokemons[pokemonDetails.Name] = pokemonDetails
			fmt.Println(fmt.Sprintf("Success ! %s caught !", pokemonDetails.Name))
			return nil
		}
		fmt.Println(fmt.Sprintf("%s Escaped !", pokemonDetails.Name))
		return nil
	}

	endpoint := fmt.Sprintf("/pokemon/%s", pokemon[0])
	pokemonDetails, err := networking.GetPokemon(endpoint)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %s", pokemonDetails.Name))
	tryCatch := calculateCatchProbability(int32(pokemonDetails.BaseExperience))
	if tryCatch {
		caughtPokemons[pokemonDetails.Name] = pokemonDetails
		fmt.Println(fmt.Sprintf("Success ! %s caught !", pokemonDetails.Name))
		return nil
	}
	fmt.Println(fmt.Sprintf("%s Escaped !", pokemonDetails.Name))
	return nil
}

func inspectCmd(pokemon ...string) error {
	if len(pokemon) == 0 {
		fmt.Printf("please pass a Pokemon name in this format : 'catch pokemon_name' \n")
		return errors.New("No Pokemon name specified")
	}
	if entry, exists := caughtPokemons[pokemon[0]]; exists {
		pokemonDetails := entry
		fmt.Println(fmt.Sprintf("\nName : %s", pokemonDetails.Name))
		fmt.Println(fmt.Sprintf("Height: %d", pokemonDetails.Height))
		fmt.Println(fmt.Sprintf("Weight: %d", pokemonDetails.Weight))
		fmt.Println("Stats : ")
		for _, statsData := range pokemonDetails.Stats {
			fmt.Println(fmt.Sprintf("  -%s: %d", statsData.Stat.Name, statsData.BaseStat))
		}
		fmt.Println("Types : ")
		for _, typesData := range pokemonDetails.Types {
			fmt.Println(fmt.Sprintf("  -%s", typesData.Type.Name))
		}
		fmt.Println()
		return nil
	}
	fmt.Println("Pokemon not available in Pokedex, you have to catch it first !")
	return nil
}

func pokedexCmd(...string) error {
	if len(caughtPokemons) == 0 {
		fmt.Println("You have 0 Pokemons caught")
		return nil
	}
	fmt.Println(fmt.Sprintf("Pokemons caught : %d ", len(caughtPokemons)))
	for _, pokemon := range caughtPokemons {
		fmt.Println("- " + pokemon.Name)
	}
	return nil
}

// define command struct
type command struct {
	name        string
	description string
	usage       string
	function    func(...string) error //use variadic because I don't know if the next assignment part is gonig to require further func arg alteration (spoiler, it doesn't)
}
