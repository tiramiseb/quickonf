package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

func Run(g []*instructions.Group) {
	instructions.SortGroups(g)
	program := tea.NewProgram(
		newModel(g),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)
	program.Start()
	commands.Clean()
}
