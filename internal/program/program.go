package program

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/global/groups"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

func Run(g []*instructions.Group) {
	groups.Initialize(g)
	toggles.Enable("filter")
	toggles.Enable("helpIntro")
	program := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	program.Start()
	commands.Clean()
}
