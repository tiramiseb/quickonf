package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func toggle(t string) tea.Cmd {
	return func() tea.Msg {
		global.Toggles.Toggle(t)
		return nil
	}
}

func enable(t string) tea.Cmd {
	return func() tea.Msg {
		global.Toggles.Enable(t)
		return nil
	}
}

func disable(t string) tea.Cmd {
	return func() tea.Msg {
		global.Toggles.Disable(t)
		return nil
	}
}
