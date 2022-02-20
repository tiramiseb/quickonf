package messages

import tea "github.com/charmbracelet/bubbletea"

// FilterMsg means the filter status must be toggled
type FilterMsg struct{}

func Filter() tea.Msg {
	return FilterMsg{}
}

type HelpMsg struct{}

func Help() tea.Msg {
	return HelpMsg{}
}
