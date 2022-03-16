package instructions

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

func (c *Command) Run(vars Variables, signalTarget chan bool) ([]*CheckReport, bool) {
	if len(c.Arguments) != len(c.Command.Arguments) {
		return []*CheckReport{{c.Command.Name, commands.StatusError, "wrong number of arguments", nil, signalTarget, "", ""}}, false
	}
	args := make([]string, len(c.Arguments))
	for i, src := range c.Arguments {
		args[i] = vars.translateVariables(src)
	}
	result, out, apply, status, before, after := c.Command.Run(args)
	for i, tgt := range c.Targets {
		if len(result) <= i {
			break
		}
		vars.define(tgt, result[i])
	}
	return []*CheckReport{{c.Command.Name, status, out, apply, signalTarget, before, after}}, status != commands.StatusError
}

func (c *Command) Reset() {
	if c.Command.Reset != nil {
		c.Command.Reset()
	}
}
