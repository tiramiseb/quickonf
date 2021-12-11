package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/state"
)

func Run(st *state.State) {
	program := tea.NewProgram(
		newModel(st),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	program.Start()
}
