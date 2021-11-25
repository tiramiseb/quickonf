package state

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

type Group struct {
	Name      string
	Commands  []Command
	variables variables
}

func (g *Group) Run(options Options) {
	out := output.NewGroup(g.Name)
	for _, command := range g.Commands {
		if !command.Run(out, g.variables, options) {
			out.Fail()
			return
		}
	}
	out.Close()
}
