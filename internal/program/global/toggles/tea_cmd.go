package toggles

import tea "github.com/charmbracelet/bubbletea"

func ToggleCmd(key string) tea.Cmd {
	return func() tea.Msg {
		Toggle(key)
		return nil
	}
}

func ToggleEnableCmd(key string) tea.Cmd {
	return func() tea.Msg {
		Enable(key)
		return nil
	}
}

func ToggleDisableCmd(key string) tea.Cmd {
	return func() tea.Msg {
		Disable(key)
		return nil
	}
}
