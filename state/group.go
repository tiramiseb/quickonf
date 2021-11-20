package state

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

type Group struct {
	Name     string
	Commands []Command
}

func (g *Group) Run() {
	out := output.NewGroup(g.Name)
	for _, command := range g.Commands {
		if ok := command.Run(out); !ok {
			out.Fail()
			return
		}
	}
	out.Close()
}
