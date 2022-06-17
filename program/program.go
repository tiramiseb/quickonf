package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/commands"
	"github.com/tiramiseb/quickonf/instructions"
)

func Run(g *instructions.Groups) {
	program := tea.NewProgram(
		newModel(g),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	program.Start()
	commands.Clean()
}
