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
var LocationsCache pokecache.Cache
var AreasCache pokecache.Cache

func init() {
	LocationsCache = *pokecache.NewCache(20)
	AreasCache = *pokecache.NewCache(20)
}

func ExploreArea(endpoint string) (AreaDetails, error) {
	cachedArea, ok := AreasCache.Get(endpoint)
	if ok {
		fmt.Println("Area cached, fetching from cache...")
		var areaDetails AreaDetails
		marshallingErr := json.Unmarshal(cachedArea, &areaDetails)
		if marshallingErr != nil {
			return AreaDetails{}, marshallingErr
		}
		return areaDetails, nil
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get(BaseURL + endpoint)
	if err != nil {
		return AreaDetails{}, err
	}

	if response.StatusCode > 399 {
		return AreaDetails{}, errors.New("Bad Request, status code : " + strconv.Itoa(response.StatusCode))
	}
	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return AreaDetails{}, err
	}
	AreasCache.Add(endpoint, data) //add to cache here

	var areaDetails AreaDetails
	marshallingErr := json.Unmarshal(data, &areaDetails)
	if marshallingErr != nil {
		return AreaDetails{}, marshallingErr
	}
	return areaDetails, nil
}

func GetLocationAreas(endpoint string) (LocationArea, error) {
	cachedLocation, ok := LocationsCache.Get(endpoint)
	if ok {
		fmt.Println("Location list cached, fetching from cache...")
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
	LocationsCache.Add(endpoint, data) //add to cache here

	//unmarshall needs a pointer to the struct, so I have to create an "instance" of it in order to pass it, GPT says it's because Unmarshal modifies the data in place instead of returning a copy, it's like that for performance reasons
	//these comments are for me to note and understand how some stuff works in go, they help me remember better than writing somewhere else
	var location LocationArea
	marshallingErr := json.Unmarshal(data, &location)
	if marshallingErr != nil {
		return LocationArea{}, marshallingErr
	}
	return location, nil
}

func GetPokemon(endpoint string) (Pokemon, error) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	response, err := client.Get(BaseURL + endpoint)
	if err != nil {
		return Pokemon{}, err
	}

	if response.StatusCode > 399 {
		return Pokemon{}, errors.New("Bad Request, status code : " + strconv.Itoa(response.StatusCode))
	}
	data, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return Pokemon{}, err
	}
	var pokemon Pokemon
	marshallingErr := json.Unmarshal(data, &pokemon)
	if marshallingErr != nil {
		return Pokemon{}, marshallingErr
	}
	return pokemon, nil
}
