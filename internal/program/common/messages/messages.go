package messages

import tea "github.com/charmbracelet/bubbletea"

// ToggleMsg means the current element should be toggled
type ToggleMsg struct{}

func Toggle() tea.Msg {
	return ToggleMsg{}
}

// HelpMsg means the help must be displayed or not
type HelpMsg struct {
	On bool
}

func Help(on bool) tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{on}
	}
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
