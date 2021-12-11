package state

import (
	"github.com/tiramiseb/quickonf/internal/commands"
)

type Command struct {
	Command   commands.Command
	Arguments []string
	Targets   []string
}

func (c *Command) Name() string {
	return c.Command.Name
}

func (c *Command) Run(out Output, vars Variables, options Options) bool {
	slow(options)
	if len(c.Arguments) != len(c.Command.Arguments) {
		out.Error("wrong number of arguments")
		return false
	}
	if len(c.Targets) > len(c.Command.Outputs) {
		out.Error("too many targets")
		return false
	}
	args := make([]string, len(c.Arguments))
	for i, src := range c.Arguments {
		args[i] = vars.translateVariables(src)
	}
	// result, ok := i.Instruction.Run(args, out, options.DryRun)
	result, ok := c.Command.Run(args, out, options.DryRun)
	if !ok {
		return false
	}
	for i, tgt := range c.Targets {
		if len(result) <= i {
			break
		}
		vars.define(tgt, result[i])
	}
	return true
}
