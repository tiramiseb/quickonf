package instructions

import (
	"github.com/tiramiseb/quickonf/commands"
)

type Command struct {
	Command   commands.Command
	Arguments []string
	Targets   []string
}

func (c *Command) Name() string {
	return c.Command.Name
}

func (c *Command) RunCheck(vars Variables, signalTarget chan bool) ([]*CheckReport, bool) {
	if len(c.Arguments) != len(c.Command.Arguments) {
		return []*CheckReport{{
			Name:         c.Command.Name,
			status:       commands.StatusError,
			message:      "wrong number of arguments",
			signalTarget: signalTarget,
		}}, false
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
	if signalTarget != nil {
		defer func() {
			signalTarget <- true
		}()
	}
	return []*CheckReport{{
		Name:         c.Command.Name,
		status:       status,
		message:      out,
		apply:        apply,
		signalTarget: signalTarget,
		Before:       before,
		After:        after,
	}}, status != commands.StatusError
}

func (c *Command) Reset() {
	if c.Command.Reset != nil {
		c.Command.Reset()
	}
}
