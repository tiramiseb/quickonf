package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	View  string
	Width int

	action tea.Cmd
}

func NewButton(text string, hintPosition int, action tea.Cmd) *Button {
	return &Button{
		View:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		Width: len(text) + 2,

		action: action,
	}
}

func (b *Button) Click() tea.Cmd {
	return b.action
}
