package instructions

import (
	"sort"
	"strings"

	"github.com/tiramiseb/quickonf/internal/output"
)

// Run is the runner for an instruction.
//
// args are the instruction arguments
// out is where the instruction must write its output
// dry is a boolean indicating the dry-run mode
// result is the instruction output values
// success is a boolean indicating if the instruction succeeded or not
type run func(args []string, out *output.Instruction, dry bool) (result []string, success bool)

// Instruction is a single instruction definition
type Instruction struct {
	Name      string   // The name of the instruction (used as a command)
	Action    string   // [used for doc] Action description
	DryRun    string   // [used for doc] Action description in dry run mode
	Arguments []string // [used for doc & counting args] Arguments description
	Outputs   []string // [used for doc & counting outputs] Outputs description
	Example   string   // [used for doc] Example(s)
	Run       run      // The function to run the instruction
}

var registry = map[string]Instruction{}

func register(instr Instruction) {
	registry[instr.Name] = instr
}

// Get returns the named instruction and a boolean, which is false if the
// instruction does not exist.
func Get(name string) (Instruction, bool) {
	instr, ok := registry[name]
	return instr, ok

}

// GetAll returns all registered instructions, sorted alphabetically
func GetAll() []Instruction {
	all := make([]Instruction, len(registry))
	i := 0
	for _, ins := range registry {
		all[i] = ins
		i++
	}
	sort.Slice(all, func(i, j int) bool {
		return strings.Compare(all[i].Name, all[j].Name) <= 0
	})
	return all
}
