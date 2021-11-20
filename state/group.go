package state

import (
	"github.com/tiramiseb/quickonf/internal/output"
)

type Group struct {
	Name     string
	Commands []*Command

	variables map[string]string
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

func (g *Group) setVariable(name, value string) {
	g.variables[name] = value
}

func (g *Group) getVariable(name string) string {
	return g.variables[name]
}
