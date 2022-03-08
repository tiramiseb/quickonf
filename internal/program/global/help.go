package global

import tea "github.com/charmbracelet/bubbletea"

func ToggleHelp() tea.Msg {
	return ToggleHelpMsg{}
}

type ToggleHelpMsg struct{}
