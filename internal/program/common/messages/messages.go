package messages

import tea "github.com/charmbracelet/bubbletea"

// ToggleMsg means the current element should be toggled
type ToggleMsg struct{}

func Toggle() tea.Msg {
	return ToggleMsg{}
}

// HelpMsg means the help must be displayed
type HelpMsg struct{}

func Help() tea.Msg {
	return HelpMsg{}
}

// FilterMsg means the filter status must be toggled
type FilterMsg struct {
	On bool
}

func Filter(on bool) tea.Cmd {
	return func() tea.Msg {
		return FilterMsg{on}
	}
}
