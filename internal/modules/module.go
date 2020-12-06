package modules

import "github.com/tiramiseb/quickonf/internal/output"

// Instruction is an instruction. Returns true if succeeds, false if there has been an error
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
	return registry[name]
}
