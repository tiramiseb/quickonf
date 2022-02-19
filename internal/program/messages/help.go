package messages

import tea "github.com/charmbracelet/bubbletea"

// HelpMsg means the help must be displayed or not
type HelpMsg struct {
	On bool
}

func Help(on bool) tea.Cmd {
	return func() tea.Msg {
		return HelpMsg{on}
	}
}
