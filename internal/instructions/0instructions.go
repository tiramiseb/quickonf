package instructions

import (
	"sort"
	"strings"

	"github.com/tiramiseb/quickonf/internal/output"
)

// Dryrun indicates if instructions must run in "dry-run" mode or not
var Dryrun = false

// Run is the runnur for an instruction
type Run func(args []string, out *output.Instruction) (result []string, success bool)

// Instruction is a single instruction definition
type Instruction struct {
	Name      string   // The name of the instruction (used as a command)
	Run       Run      // The function to run the instruction
	Action    string   // [used for doc] Action description
	DryRun    string   // [used for doc] Action description in dry run mode
	Arguments []string // [used for doc & counting args] Arguments description
	Outputs   []string // [used for doc & counting outputs] Outputs description
	Example   string   // [used for doc] Example(s)
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
