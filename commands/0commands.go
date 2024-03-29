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
	RunningPercent(message string) chan int
	RunningPercentf(format string, a ...interface{}) chan int
	Success(message string)
	Successf(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
}

type Status int

const (
	StatusNotRun Status = iota
	StatusInfo
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
// status is the status after checking
// before is the value before modification (only if apply is not nil)
// after is the value after modification (only if apply is not nil)
type run func(args []string) (result []string, message string, apply Apply, status Status, before, after string)

// Command is a single command definition
type Command struct {
	Name      string   `json:"name"`       // The name of the command
	Action    string   `json:"action"`     // [used for doc] Action description
	Arguments []string `json:"arguments"`  // [used for doc & counting args] Arguments description
	Outputs   []string `json:"outputs"`    // [used for doc & counting outputs] Outputs description
	Example   string   `json:"example"`    // [used for doc] Example(s)
	Run       run      `yaml:"-" json:"-"` // The function to run the command
	Reset     func()   `yaml:"-" json:"-"` // The function to reset data in order to cleanly re-run the command
}

// Apply is the process to apply a given instance of a command on the system
type Apply func(Output) (success bool)

var registry = map[string]*Command{}

func register(cmd *Command) {
	registry[cmd.Name] = cmd
}

// Get returns the named command and a boolean, which is false if the
// command does not exist.
func Get(name string) (*Command, bool) {
	instr, ok := registry[name]
	return instr, ok
}

// UGet returns the named command
func UGet(name string) *Command {
	return registry[name]
}

// GetAll returns all registered commands, sorted alphabetically
func GetAll() []*Command {
	all := make([]*Command, len(registry))
	i := 0
	for _, cmd := range registry {
		all[i] = cmd
		i++
	}
	sort.Slice(all, func(i, j int) bool {
		return strings.Compare(all[i].Name, all[j].Name) <= 0
	})
	return all
}

func ListStartWith(prefix string) []*Command {
	var commands []*Command
	for _, cmd := range registry {
		if strings.HasPrefix(cmd.Name, prefix) {
			commands = append(commands, cmd)
		}
	}
	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].Name, commands[j].Name) <= 0
	})
	return commands
}

func (c *Command) LastArgumentIsVariadic() bool {
	return len(c.Arguments) > 0 && strings.HasSuffix(c.Arguments[len(c.Arguments)-1], "...")
}

func (c *Command) LastOutputIsVariadic() bool {
	return len(c.Outputs) > 0 && strings.HasSuffix(c.Outputs[len(c.Outputs)-1], "...")
}

func (c *Command) MarkdownHelp() string {
	var result strings.Builder

	result.WriteString("**")
	result.WriteString(c.Name)
	result.WriteString("**: ")
	result.WriteString(c.Action)
	result.WriteString("\n\n")

	if len(c.Arguments) > 0 {
		result.WriteString("Arguments:\n\n")
		for _, arg := range c.Arguments {
			result.WriteString("- ")
			result.WriteString(arg)
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}

	if len(c.Outputs) > 0 {
		result.WriteString("Outputs:\n\n")
		for _, out := range c.Outputs {
			result.WriteString("- ")
			result.WriteString(out)
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}

	if c.Example != "" {
		result.WriteString("Example:\n\n```\n")
		result.WriteString(c.Example)
		result.WriteString("\n```\n")
	}

	return result.String()
}
