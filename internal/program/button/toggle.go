package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Toggle struct {
	status bool

	offView string
	onView  string
	Width   int

	Click tea.Cmd
}

func NewToggle(text string, hintPosition int, action tea.Cmd, initial bool) *Toggle {
	return &Toggle{
		status: initial,

		offView: lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		onView:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, activeHintStyle, activeStyle),
		Width:   len(text) + 2,

		Click: action,
	}
}

func (t *Toggle) View() string {
	if t.status {
		return t.onView
	}
	return t.offView
}

func (t *Toggle) ChangeStatus(status bool) *Toggle {
	t.status = status
	return t
}
