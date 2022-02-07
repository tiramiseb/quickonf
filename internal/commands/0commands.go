package commands

import (
	"sort"
	"strings"
)

type Output interface {
	Info(message string)
	Infof(format string, a ...interface{})
	Running(message string)
	Runningf(format string, a ...interface{})
	Success(message string)
	Successf(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
}

type Status int

const (
	StatusInfo Status = iota
	StatusRunning
	StatusSuccess
	StatusError
)

// apply is a generated function that applies the result of the command to the system. It is returned by the run command, for potential future call.
//
// Output is written to the output given to the run command
//
// success is a boolean indicating if the command succeeded or not
// type Apply func() (success bool)

// run is the function that runs the command and prepares an apply function
//
// args are the command arguments (same as for "run")
// result is the command output values
// message is the message to display to the user
// apply is the function that would be executed in order to apply the command to the system. If nil, it means there is nothing to apply (the system is already in the requested state)
// success is a boolean indicating if preparing the apply function succeeded or not
type run func(args []string) (result []string, message string, apply *Apply, status Status)

// Command is a single command definition
type Command struct {
	Name      string   // The name of the command
	Action    string   // [used for doc] Action description
	Arguments []string // [used for doc & counting args] Arguments description
	Outputs   []string // [used for doc & counting outputs] Outputs description
	Example   string   // [used for doc] Example(s)
	Run       run      // The function to run the command
	Reset     func()   // The function to reset data in order to cleanly re-run the command
}

// Apply is the process to apply a given instance of a command on the system
type Apply struct {
	Name  string
	Intro string
	Run   func(Output) (success bool)
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
