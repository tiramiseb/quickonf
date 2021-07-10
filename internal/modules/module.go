package modules

import (
	"errors"

	"github.com/tiramiseb/quickonf/internal/output"
)

// Instruction is an instruction. Returns nil if succeeds, or an error
type Instruction func(in interface{}, out output.Output) error

var registry = map[string]Instruction{}

// Dryrun allows running instructions without system modification
var Dryrun = false

// Register adds an instruction
func Register(name string, instruction Instruction) {
	registry[name] = instruction
}

// Get gets an instruction
func Get(name string) Instruction {
	instruction, ok := registry[name]
	if !ok {
		return func(interface{}, output.Output) error {
			return errors.New("[No instruction named \"" + name + "\"]")
		}
	}
	return instruction
}
