package main

import (
	"bufio"
	"fmt"
	"os"
)

type command struct {
	name        string
	description string
	function    func() error
}

func fetchData(pokemon string) (string, error) {
	mockdata := "just some random stuff, this should be structured JSON data then formatted"
	return mockdata, nil
}

// todo : create fetching function, create help command and structs
func main() {
	exitFlag := false
	fmt.Println("Welcome to the GoGo second generation PokeDex")
	fmt.Println("For help, please type /help")
	reader := bufio.NewScanner(os.Stdin)

	for !exitFlag {
		for reader.Scan() {
			text := reader.Text()
			//switch that tests for commands
			switch text {
			case "/help":
				fmt.Println("print help here")

			case "/exit":
				fmt.Println("Exiting. . .")
				exitFlag = true

			default:
				returnedData, fetchErr := fetchData(text)
				if fetchErr != nil {
					fmt.Println(returnedData)
				} else {
					fmt.Println(fetchErr)
				}
			}
			if exitFlag {
				break
			}
		}
	}

}
