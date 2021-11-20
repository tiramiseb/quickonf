package instructions

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

// Dryrun indicates if instructions must run in "dry-run" mode or not
var Dryrun = false

// Instruction is a quickonf instruction
type Instruction func(args []string, out *output.Instruction) (result []string, success bool)

var (
	registry       = map[string]Instruction{}
	nbArguments    = map[string]int{}
	nbOutputValues = map[string]int{}
)

func register(name string, instr Instruction, arguments int, outputValues int) {
	registry[name] = instr
	nbArguments[name] = arguments
	nbOutputValues[name] = outputValues
}

// Get returns the named instruction, the number of arguments it needs,
// the number of values it returns and a boolean, which is false if the
// instruction does not exist.
func Get(name string) (Instruction, int, int, bool) {
	instr, ok := registry[name]
	args := nbArguments[name]
	out := nbOutputValues[name]
	return instr, args, out, ok

}
