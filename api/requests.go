package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	PokemonEndpoint       string = "https://pokeapi.co/api/v2/pokemon/"
	LocationAreaEndpoint  string = "https://pokeapi.co/api/v2/location-area/"
	LocationAreaFirstPage string = LocationAreaEndpoint + "?offset=0&limit=20"
)

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

// GetPokemonDetails polls the pokeapi for details on the given pokemon.
func GetPokemonDetails(pokemonName string) PokemonDetails {
	body := pollApi(PokemonEndpoint + pokemonName)
	var details PokemonDetails
	err := json.Unmarshal(body, &details)
	if err != nil {
		log.Fatal(err)
	}
	return details
}
