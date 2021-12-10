package commands

import (
	"sort"
	"strings"
)

type output interface {
	Info(message string)
	Infof(format string, a ...interface{})
	Success(message string)
	Successf(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
}

// Run is the runner for a command.
//
// args are the command arguments
// out is where the command must write its output
// dry is a boolean indicating the dry-run mode
// result is the command output values
// success is a boolean indicating if the command succeeded or not
type run func(args []string, out output, dry bool) (result []string, success bool)

// Command is a single command definition
type Command struct {
	Name      string   // The name of the command
	Action    string   // [used for doc] Action description
	DryRun    string   // [used for doc] Action description in dry run mode
	Arguments []string // [used for doc & counting args] Arguments description
	Outputs   []string // [used for doc & counting outputs] Outputs description
	Example   string   // [used for doc] Example(s)
	Run       run      // The function to run the command
}

var registry = map[string]Command{}

func register(cmd Command) {
	registry[cmd.Name] = cmd
}

// Get returns the named command and a boolean, which is false if the
// command does not exist.
func Get(name string) (Command, bool) {
	instr, ok := registry[name]
	return instr, ok

}

// GetAll returns all registered commands, sorted alphabetically
func GetAll() []Command {
	all := make([]Command, len(registry))
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
