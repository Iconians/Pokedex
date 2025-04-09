package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
}

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func MapCommand(config *Config) {
	url := config.Next
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to fetch location areas:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Failed to parse response:", err)
		return
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous
}

func MapBackCommand(config *Config) {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return
	}

	res, err := http.Get(config.Previous)
	if err != nil {
		fmt.Println("Failed to fetch previous locations:", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	var data LocationAreaResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Failed to parse response:", err)
		return
	}

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	config.Next = data.Next
	config.Previous = data.Previous
}
