package instructions

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

// Dryrun indicates if instructions must run in "dry-run" mode or not
var Dryrun = false

// Run is the runnur for an instruction
type Run func(args []string, out *output.Instruction) (result []string, success bool)

// Instruction is a single instruction definition
type Instruction struct {
	Name            string
	Run             Run
	NumberArguments int
	NumberOutputs   int
}

var (
	registry = map[string]Instruction{}
)

func register(instr Instruction) {
	registry[instr.Name] = instr
}

// Get returns the named instruction and a boolean, which is false if the
// instruction does not exist.
func Get(name string) (Instruction, bool) {
	instr, ok := registry[name]
	return instr, ok

}
