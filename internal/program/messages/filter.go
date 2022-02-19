package messages

import tea "github.com/charmbracelet/bubbletea"

// FilterMsg means the filter status must be toggled
type FilterMsg struct {
	Enable bool
}

func Filter(enable bool) tea.Cmd {
	return func() tea.Msg {
		return FilterMsg{enable}
	}
}
