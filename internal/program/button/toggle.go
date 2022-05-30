package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/messages"
)

type Toggle struct {
	key    string
	status bool

	offView string
	onView  string
	Width   int
}

func NewToggle(text string, hintPosition int, key string, initial bool) *Toggle {
	return &Toggle{
		key:    key,
		status: initial,

		offView: lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		onView:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, activeHintStyle, activeStyle),
		Width:   len(text) + 2,
	}
}

func (t *Toggle) Click() tea.Msg {
	return messages.Toggle{Name: t.key, Action: messages.ToggleActionToggle}
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
