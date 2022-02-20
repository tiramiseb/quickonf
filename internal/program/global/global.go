package global

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	BgColor = lipgloss.Color(fmt.Sprintf("%v", termenv.BackgroundColor()))
	FgColor = lipgloss.Color(fmt.Sprintf("%v", termenv.ForegroundColor()))
)

type global struct {
	values map[string]bool
}

var Global = &global{
	values: map[string]bool{},
}

func (g *global) Set(key string, value bool) {
	g.values[key] = value
}

func (g *global) Get(key string) bool {
	return g.values[key]
}
