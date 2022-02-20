package messages

import tea "github.com/charmbracelet/bubbletea"

type DetailsMsg struct{}

func Details() tea.Msg {
	return DetailsMsg{}
}

type FilterMsg struct{}

func Filter() tea.Msg {
	return FilterMsg{}
}

type HelpMsg struct{}

func Help() tea.Msg {
	return HelpMsg{}
}
