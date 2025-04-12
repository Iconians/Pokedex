package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Stat struct {
	BaseStat int `json:"base_stat"`
	StatInfo struct {
		Name string `json:"name"`
	} `json:"stat"`
}

type TypeInfo struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type Pokemon struct {
	Name           string     `json:"name"`
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []Stat     `json:"stats"`
	Types          []TypeInfo `json:"types"`
}

func GetPokemon(cfg *Config, name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	if data, ok := cfg.Cache.Get(url); ok {
		var p Pokemon
		if err := json.Unmarshal(data, &p); err != nil {
			return Pokemon{}, err
		}
		return p, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return Pokemon{}, fmt.Errorf("could not find pokemon %s", name)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var p Pokemon
	if err := json.Unmarshal(body, &p); err != nil {
		return Pokemon{}, err
	}

	cfg.Cache.Add(url, body)
	return p, nil
}
