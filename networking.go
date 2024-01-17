package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

var baseURL string = "https://pokeapi.co/api/v2"
var client http.Client = http.Client{
	Timeout: 30 * time.Second,
}

func getLocationAreas(client http.Client) (LocationArea, error) {
	endpoint := "/location-area"
	response, err := client.Get(baseURL + endpoint)
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
	//unmarshall needs a pointer to the struct, so I have to create an "instance" of it in order to pass it, GPT says it's because Unmarshal modifies the data in place instead of returning a copy, it's like that for performance reasons

	var locationArea LocationArea
	marshallingErr := json.Unmarshal(data, &locationArea)
	if marshallingErr != nil {
		return LocationArea{}, marshallingErr
	}

	return locationArea, nil

}

type LocationArea struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
