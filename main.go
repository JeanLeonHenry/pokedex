package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	urls "net/url"
	"os"
	"strconv"
	"time"

	"github.com/JeanLeonHenry/pokedex/api"
	"github.com/JeanLeonHenry/pokedex/pokecache"
)

type command struct {
	name        string
	description string
	fn          func()
}

func (c command) String() string {
	return fmt.Sprintf("%v: %v", c.name, c.description)
}

var cmds map[string]command

func displayHelp() {
	fmt.Println(`Pokedex

Usage: pokedex <command>

Commands:`)
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

func (c *config) PrintLocations(url string) {
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

func (c *config) Next() {
	c.PrintLocations(c.next)
}
func (c *config) Prev() {
	c.PrintLocations(c.previous)
}

func main() {
	// Set up
	cfg := &config{next: api.BaseUrl, previous: api.BaseUrl, cache: *pokecache.NewCache(20 * time.Second)}
	cmds = map[string]command{
		"map":  {name: "map", description: "Display next 20 locations.", fn: cfg.Next},
		"mapb": {name: "mapb", description: "Display previous 20 locations.", fn: cfg.Prev},
		"help": {name: "help", description: "Display help message.", fn: displayHelp},
		"exit": {name: "exit", description: "Quit program.", fn: func() { os.Exit(0) }},
	}
	// REPL
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		if ok := scanner.Scan(); !ok {
			log.Fatal("Wrong input. Quitting.")
		}
		input := scanner.Text()
		if _, ok := cmds[input]; !ok {
			log.Println("Wrong cmd:", input)
			continue
		}
		cmds[input].fn()
	}
}
