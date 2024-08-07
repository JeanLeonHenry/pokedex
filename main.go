package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/JeanLeonHenry/pokedex/api"
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
}

func (c *config) PrintLocations(url string) {
	if url == "" {
		log.Println("Empty url")
		return
	}
	locations, previous, next := api.GetLocationsPage(url)
	c.previous = previous
	c.next = next
	fmt.Println(locations)
}

func (c *config) Next() {
	c.PrintLocations(c.next)
}
func (c *config) Prev() {
	c.PrintLocations(c.previous)
}

func main() {
	// Set up
	cfg := &config{next: api.BaseUrl, previous: api.BaseUrl}
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
