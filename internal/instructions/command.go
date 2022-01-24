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

func (c *Command) Run(vars Variables) ([]commands.Apply, []CheckReport, bool) {
	if len(c.Arguments) != len(c.Command.Arguments) {
		return nil, []CheckReport{{c.Command.Name, commands.StatusError, "wrong number of arguments"}}, false
	}
	args := make([]string, len(c.Arguments))
	for i, src := range c.Arguments {
		args[i] = vars.translateVariables(src)
	}
	result, out, apply, status := c.Command.Run(args)
	for i, tgt := range c.Targets {
		if len(result) <= i {
			break
		}
		vars.define(tgt, result[i])
	}
	var applies []commands.Apply
	if apply != nil {
		applies = append(applies, *apply)
	}
	return applies, []CheckReport{{c.Command.Name, status, out}}, status != commands.StatusError
}

func (c *Command) Reset() {
	if c.Command.Reset != nil {
		c.Command.Reset()
	}
}
