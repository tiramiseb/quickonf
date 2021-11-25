package state

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

type If struct {
	Operation Operation
	Commands  []Command
}

func (i *If) Run(groupOut *output.Group, vars variables, options Options) bool {
	success := i.Operation.Compare(vars)
	if !success {
		out := groupOut.NewInstruction("=")
		out.Infof(`"%s" is false, not running commands...`, i.Operation.String())
		return true
	}
	slow(options)
	for _, cmd := range i.Commands {
		if !cmd.Run(groupOut, vars, options) {
			return false
		}
	}
	return true
}
