package main

import (
	"fmt"
	"os"

	"github.com/stygian91/veggies/commands"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Not enough arguments.")
		return
	}

	rest := args[1:]

	cmdname := args[0]
	switch cmdname {
	case "new":
		if len(rest) == 0 {
			fmt.Println("Not enough arguments for the 'new' command.")
			return
		}
		if err := commands.New(rest[0]); err != nil {
			panic(err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmdname)
	}
}
