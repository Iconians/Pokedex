package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Encounter struct {
	Pokemon struct {
		Name string `json:"name"`
	} `json:"pokemon"`
}

const baseURL = "https://pokeapi.co/api/v2"

func ExploreArea(cfg *Config, areaName string) error {
	url := baseURL + "/location-area/" + areaName

	data, ok := cfg.Cache.Get(url)
	if ok {
		return printPokemonInArea(data, areaName)
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to fetch location area: %s", areaName)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	cfg.Cache.Add(url, body)
	return printPokemonInArea(body, areaName)
}

func printPokemonInArea(data []byte, areaName string) error {
	var parsed struct {
		PokemonEncounters []Encounter `json:"pokemon_encounters"`
	}

	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Println("Found Pokemon:")

	for _, e := range parsed.PokemonEncounters {
		fmt.Printf(" - %s\n", e.Pokemon.Name)
	}
	return nil
}
