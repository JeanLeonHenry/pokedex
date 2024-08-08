package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	LocationAreaEndpoint  string = "https://pokeapi.co/api/v2/location-area/"
	LocationAreaFirstPage string = LocationAreaEndpoint + "?offset=0&limit=20"
)

type LocationAreaResponse struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  LocationSlice `json:"results"`
}
type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             Location              `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

func (l LocationArea) String() string { return l.Name }

type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterVersionDetail struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}
type EncounterMethodRate struct {
	EncounterMethod EncounterMethod          `json:"encounter_method"`
	VersionDetails  []EncounterVersionDetail `json:"version_details"`
}
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type LocationSlice []Location

func (l LocationSlice) String() (res string) {
	for _, location := range l {
		res += location.Name + "\n"
	}
	return
}

type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Name struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}
type PokemonSlice []Pokemon

func (p PokemonSlice) String() (res string) {
	for _, pokemon := range p {
		res += fmt.Sprintln("\t-", pokemon.Name)
	}
	return
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Method struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type EncounterDetails struct {
	MinLevel        int    `json:"min_level"`
	MaxLevel        int    `json:"max_level"`
	ConditionValues []any  `json:"condition_values"`
	Chance          int    `json:"chance"`
	Method          Method `json:"method"`
}
type VersionEncounterDetail struct {
	Version          Version            `json:"version"`
	MaxChance        int                `json:"max_chance"`
	EncounterDetails []EncounterDetails `json:"encounter_details"`
}
type PokemonEncounter struct {
	Pokemon        Pokemon                  `json:"pokemon"`
	VersionDetails []VersionEncounterDetail `json:"version_details"`
}

func pollApi(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return body
}

// GetLocationsPage polls the pokeapi for an api.Limit number of location areas, starting from given page.
func GetLocationsPage(url string) LocationAreaResponse {
	if url == "" {
		url = LocationAreaFirstPage
	}
	body := pollApi(url)
	var locations LocationAreaResponse
	err := json.Unmarshal(body, &locations)
	if err != nil {
		log.Fatal(err)
	}
	return locations
}

// GetPokemonsInArea polls the pokeapi for the given location and returns the local pokemons.
func GetPokemonsInArea(locationName string) (result PokemonSlice) {
	body := pollApi(LocationAreaEndpoint + locationName)
	var location LocationArea
	err := json.Unmarshal(body, &location)
	if err != nil {
		log.Fatal(err)
	}
	for _, encounter := range location.PokemonEncounters {
		result = append(result, encounter.Pokemon)
	}
	return result
}
