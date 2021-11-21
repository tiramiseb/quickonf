package state

import (
	"time"

	"github.com/tiramiseb/quickonf/internal/output"
)

type If struct {
	Operation Operation
	Commands  []Command
}

func (i *If) Run(groupOut *output.Group, vars variables) bool {
	success := i.Operation.Compare(vars)
	if !success {
		out := groupOut.NewInstruction("=")
		out.Infof(`"%s" is false, not running commands...`, i.Operation.String())
		time.Sleep(2 * time.Second)
		return true
	}
	for _, cmd := range i.Commands {
		if !cmd.Run(groupOut, vars) {
			return false
		}
	}
	return true
}
