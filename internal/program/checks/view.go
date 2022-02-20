package checks

import (
	"strings"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) RedrawView() *Model {
	filter := global.Global.Get("filter")
	var view string
	for _, g := range m.groups {
		status := g.Status()
		if filter && status == commands.StatusSuccess {
			continue
		}
		line := " " + g.Name
		remaining := m.width - len(line)
		switch {
		case remaining < 0:
			line = line[:m.width-1]
		case remaining > 0:
			line += strings.Repeat(" ", remaining)
		}
		view += styles[status].Render(line) + "\n"

	}
	m.completeView = view
	return m
}

func (m *Model) View() string {
	return m.completeView
}
