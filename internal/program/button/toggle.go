package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

type Toggle struct {
	key string

	offView string
	onView  string
	Width   int
}

func NewToggle(text string, hintPosition int, key string) *Toggle {
	return &Toggle{
		key: key,

		offView: lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		onView:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, activeHintStyle, activeStyle),
		Width:   len(text) + 2,
	}
}

func (t *Toggle) Click() tea.Cmd {
	return toggles.ToggleCmd(t.key)
}

func (t *Toggle) View() string {
	if toggles.Get(t.key) {
		return t.onView
	}
	return t.offView
}
