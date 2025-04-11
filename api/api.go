package api

import (
	"encoding/json"
	"fmt"
	"iconians/pokedexcli/internals/pokecache"
	"io"
	"net/http"
)

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
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

func fetchURL(url string, config *Config) ([]byte, error) {
	if config.Cache != nil {
		if val, found := config.Cache.Get(url); found {
			fmt.Println("Using Cached data.")
			return val, nil
		}
	}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if config.Cache != nil {
		config.Cache.Add(url, body)
	}
	return body, nil
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

	body, err := fetchURL(config.Previous, config)
	// res, err := http.Get(config.Previous)
	if err != nil {
		fmt.Println("Failed to fetch previous locations:", err)
		return
	}
	// defer res.Body.Close()

	// body, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println("Failed to read response:", err)
	// 	return
	// }

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
