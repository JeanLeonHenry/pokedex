package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

func main() {
	cmds = map[string]command{
		"map":  {name: "map", description: "Display next 20 locations.", fn: notImplemented},
		"mapb": {name: "mapb", description: "Display previous 20 locations.", fn: notImplemented},
		"help": {name: "help", description: "Display help message.", fn: displayHelp},
		"exit": {name: "exit", description: "Quit program.", fn: func() { os.Exit(0) }},
	}
	// REPL
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("pokedex > ")
		if ok := scanner.Scan(); !ok {
			log.Println("Wrong input")
		}
		input := scanner.Text()
		if _, ok := cmds[input]; !ok {
			log.Println("Wrong cmd:", input)
			continue
		}
		cmds[input].fn()
	}
}
