package modules

import (
	"fmt"

	"github.com/tiramiseb/quickonf/internal/output"
)

// Instruction is an instruction. Returns nil if succeeds, or an error
type Instruction func(in interface{}, out output.Output) error

var registry = map[string]Instruction{}

// Dryrun allows running instructions without system modification
var Dryrun = false

// Basepath is the base path for calculating files paths (based on configuration paths)
var Basepath = ""

// Register adds an instruction
func Register(name string, instruction Instruction) {
	registry[name] = instruction
}

// Get gets an instruction
func Get(name string) Instruction {
	instruction, ok := registry[name]
	if !ok {
		return func(interface{}, output.Output) error {
			return fmt.Errorf(`[no instruction named "%s"]`, name)
		}
	}
	return instruction
}
