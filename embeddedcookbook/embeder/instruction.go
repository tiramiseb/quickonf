package main

import (
	"fmt"

	"github.com/tiramiseb/quickonf/instructions"
)

func (e *embeder) instruction(level int, instr instructions.Instruction) {
	switch instr := instr.(type) {
	case *instructions.Command:
		e.command(level, instr)
	case *instructions.Expand:
		e.expand(level, instr)
	case *instructions.If:
		e.ifthen(level, instr)
	case *instructions.Recipe:
		e.recipe(level, instr)
	default:
		panic(fmt.Sprintf("Unknown instruction: %#v", instr))
	}
}
