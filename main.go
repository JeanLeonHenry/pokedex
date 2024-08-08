package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	urls "net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JeanLeonHenry/pokedex/api"
	"github.com/JeanLeonHenry/pokedex/pokecache"
)

type command struct {
	name        string
	description string
	fn          func(...string)
}

func (c command) String() string {
	return fmt.Sprintf("%v: %v", c.name, c.description)
}

var cmds map[string]command

func displayHelp(...string) {
	fmt.Println(`Pokedex

Usage: pokedex <command>

Commands:`)
	// BUG: map traversal order isn't deterministic
	for _, cmd := range cmds {
		fmt.Println("\t", cmd)
	}
}

func notImplemented() {
	fmt.Println("Not implemented")
}

type config struct {
	next     string
	previous string
	cache    pokecache.Cache
}

func (c *config) printLocations(url string) {
	if url == "" {
		fmt.Println("Can't go back from first page.")
		return
	}
	var response api.LocationAreaResponse
	if data, ok := c.cache.Get(url); !ok {
		response = api.GetLocationsPage(url)
		c.previous = response.Previous
		c.next = response.Next

		dataToCache, err := json.Marshal(response)
		if err != nil {
			log.Fatalln("Error: couldn't cache response for ", url)
		}
		c.cache.Add(url, dataToCache)
	} else {
		err := json.Unmarshal(data, &response)
		if err != nil {
			log.Fatalln("Error: couldn't unpack cache entry for ", url)
		}
		c.previous = response.Previous
		c.next = response.Next
	}
	fmt.Println(response.Results)
	current, _ := urls.Parse(url)
	offset, _ := strconv.Atoi(current.Query()["offset"][0])
	fmt.Println("Results from", offset, "to", offset+19)
}

func (c *config) printPokemons(args ...string) {
	if len(args) != 1 {
		log.Println("explore usage: explore <location name>")
		return
	}
	locationName := args[0]

	// NOTE: implement a little spinner that cycles through . -> .. -> ... ?
	fmt.Println("Exploring", locationName, "...")

	var pokemons api.PokemonSlice
	url := api.LocationAreaEndpoint + locationName
	if data, ok := c.cache.Get(url); !ok {
		pokemons = api.GetPokemonsInArea(locationName)
		dataToCache, err := json.Marshal(pokemons)
		if err != nil {
			log.Fatalln("Error: couldn't cache response for ", url)
		}
		c.cache.Add(url, dataToCache)
	} else {
		err := json.Unmarshal(data, &pokemons)
		if err != nil {
			log.Fatalln("Error: couldn't unpack cache entry for ", url)
		}
	}
	fmt.Println("Found Pokemon:")
	fmt.Println(pokemons)
}

func (c *config) Next(...string) {
	c.printLocations(c.next)
}
func (c *config) Prev(...string) {
	c.printLocations(c.previous)
}

func main() {
	// Set up
	cfg := &config{next: api.LocationAreaFirstPage, previous: api.LocationAreaFirstPage, cache: *pokecache.NewCache(20 * time.Second)}
	cmds = map[string]command{
		"map":     {name: "map", description: "Display next 20 locations.", fn: cfg.Next},
		"mapb":    {name: "mapb", description: "Display previous 20 locations.", fn: cfg.Prev},
		"explore": {name: "explore <location>", description: "List pokemons in the given location.", fn: cfg.printPokemons},
		"help":    {name: "help", description: "Display help message.", fn: displayHelp},
		"exit":    {name: "exit", description: "Quit program.", fn: func(...string) { os.Exit(0) }},
	}
	// REPL
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		if ok := scanner.Scan(); !ok {
			log.Fatal("Wrong input. Quitting.")
		}
		input := scanner.Text()
		args := strings.Fields(strings.TrimSpace(input))
		if len(args) == 0 {
			log.Println("Wrong command.")
			continue
		}
		if cmd, ok := cmds[args[0]]; !ok {
			log.Println("Wrong command.")
			continue
		} else {
			cmd.fn(args[1:]...)
		}
	}
}
