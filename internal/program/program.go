package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

func Run(g *instructions.Groups) {
	toggles.Enable("filter")
	toggles.Enable("helpIntro")
	program := tea.NewProgram(
		newModel(g),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	program.Start()
	commands.Clean()
}
