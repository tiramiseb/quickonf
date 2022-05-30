package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
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
