package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands"
)

func list() {
	fmt.Println("Available commandss:")
	var maxNameSize int
	all := commands.GetAll()
	for _, i := range all {
		l := len(i.Name)
		if l > maxNameSize {
			maxNameSize = l
		}
	}
	for _, i := range all {
		fmt.Printf(" * %-*s: %s\n", maxNameSize, i.Name, i.Action)
	}
}

func commandHelp(cmd string) {
	command, ok := commands.Get(cmd)
	if !ok {
		fmt.Printf(`Command "%s" does not exist`, cmd)
		os.Exit(1)
	}
	fmt.Printf("Command \"%s\" usage\n", cmd)
	fmt.Printf("      Action: %s\n", command.Action)
	fmt.Printf("     Dry-run: %s\n", command.DryRun)
	if len(command.Arguments) > 0 {
		fmt.Printf("   Arguments:\n")
		for i, args := range command.Arguments {
			fmt.Printf("       #%d - %s\n", i+1, args)

		}
	}
	if len(command.Outputs) > 0 {
		fmt.Printf("     Outputs:\n")
		for i, out := range command.Outputs {
			fmt.Printf("    #%d - %s\n", i+1, out)
		}
	}
	fmt.Printf("  Example(s):\n")
	for _, l := range strings.Split(command.Example, "\n") {
		fmt.Printf("       | %s\n", l)
	}
}
