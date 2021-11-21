package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tiramiseb/quickonf/internal/instructions"
)

func list() {
	fmt.Println("Available instructions:")
	var maxNameSize int
	all := instructions.GetAll()
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

func instructionHelp(instruction string) {
	instr, ok := instructions.Get(instruction)
	if !ok {
		fmt.Printf(`Instruction "%s" does not exist`, instruction)
		os.Exit(1)
	}
	fmt.Printf("Instruction \"%s\" usage\n", instruction)
	fmt.Printf("      Action: %s\n", instr.Action)
	fmt.Printf("     Dry-run: %s\n", instr.DryRun)
	if len(instr.Arguments) > 0 {
		fmt.Printf("   Arguments:\n")
		for i, args := range instr.Arguments {
			fmt.Printf("       #%d - %s\n", i+1, args)

		}
	}
	if len(instr.Outputs) > 0 {
		fmt.Printf("     Outputs:\n")
		for i, out := range instr.Outputs {
			fmt.Printf("    #%d - %s\n", i+1, out)
		}
	}
	fmt.Printf("  Example(s):\n")
	for _, l := range strings.Split(instr.Example, "\n") {
		fmt.Printf("       | %s\n", l)
	}
}
