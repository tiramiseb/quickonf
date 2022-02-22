package details

import (
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) ChangeView(idx int) *Model {
	m.displayedGroup = idx
	return m
}

func (m *Model) View() string {
	if m.displayedGroup < 0 {
		m.displayedGroup = 0
	} else if m.displayedGroup >= len(m.groups) {
		m.displayedGroup = len(m.groups) - 1
	}
	var view string
	if len(m.groups[m.displayedGroup].Reports) == 0 {
		for _, ins := range m.groups[m.displayedGroup].Instructions {
			view += global.Styles[commands.StatusNotRun].Render(
				global.MakeWidth(ins.Name(), m.width),
			) + "\n"
		}
	}
	for _, rep := range m.groups[m.displayedGroup].Reports {
		view += global.Styles[rep.Status].Render(
			global.MakeWidth(rep.Message, m.width),
		) + "\n"
	}
	// return view
	m.viewport.SetContent(view)
	return m.viewport.View()
}
