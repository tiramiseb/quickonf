package main

import (
	"fmt"

	"github.com/tiramiseb/quickonf/instructions"
)

func (e *embeder) ifthen(level int, instr *instructions.If) {
	e.write(level, "&If{")
	switch op := instr.Operation.(type) {
	case *instructions.Different:
		e.write(level+1, `Operation: &Different{Left: "%s", Right: "%s"},`, op.Left, op.Right)
	case *instructions.Equal:
		e.write(level+1, `Operation: &Equal{Left: "%s", Right: "%s"},`, op.Left, op.Right)
	case *instructions.FileAbsent:
		e.write(level+1, `Operation: &FileAbsent{Path: "%s"},`, op.Path)
	case *instructions.FilePresent:
		e.write(level+1, `Operation: &FilePresent{Path: "%s"},`, op.Path)
	default:
		panic(fmt.Sprintf("Unknown operation: %#v", op))
	}
	e.write(level+1, "Instructions: []Instruction{")
	for _, i := range instr.Instructions {
		e.instruction(level+2, i)
	}
	e.write(level+1, "},")
	e.write(level, "},")
}
