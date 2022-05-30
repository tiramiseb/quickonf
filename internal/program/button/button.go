package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	View  string
	Width int

	Click tea.Cmd
}

func NewButton(text string, hintPosition int, action tea.Cmd) *Button {
	return &Button{
		View:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		Width: len(text) + 2,

		Click: action,
	}
}
