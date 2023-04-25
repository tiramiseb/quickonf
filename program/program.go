package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/commands"
)

func Run(config string) {
	program := tea.NewProgram(
		newModel(config),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	program.Start()
	commands.Clean()
}
