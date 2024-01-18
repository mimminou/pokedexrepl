package networking

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mimminou/pokedexrepl/internal/pokecache"
	"io"
	"net/http"
	"strconv"
	"time"
)

var BaseURL string = "https://pokeapi.co/api/v2"
var cache pokecache.Cache

func init() {
	cache = *pokecache.NewCache(5)
}

func GetLocationAreas(endpoint string) (LocationArea, error) {
	cachedLocation, ok := cache.Get(endpoint)
	if ok {
		fmt.Println("Checking the cache")
		var location LocationArea
		marshallingErr := json.Unmarshal(cachedLocation, &location)
		if marshallingErr != nil {
			return LocationArea{}, marshallingErr
		}
		return location, nil
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get(BaseURL + endpoint)
	if err != nil {
		return LocationArea{}, err
	}

	if response.StatusCode > 399 {
		return LocationArea{}, errors.New("Bad Request, status code : " + strconv.Itoa(response.StatusCode))
	}
	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return LocationArea{}, err
	}
	cache.Add(endpoint, data) //add to cache here

	//unmarshall needs a pointer to the struct, so I have to create an "instance" of it in order to pass it, GPT says it's because Unmarshal modifies the data in place instead of returning a copy, it's like that for performance reasons
	//these comments are for me to note and understand how some stuff works in go, they help me remember better than writing somewhere else
	var location LocationArea
	marshallingErr := json.Unmarshal(data, &location)
	if marshallingErr != nil {
		return LocationArea{}, marshallingErr
	}
	return location, nil
}

// I don't like this struct being here, but I don't want to get stuff complicated by creating other packages so I'll deal with it later
type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
