package button

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/global"
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

func (t *Toggle) Click() {
	global.Toggles.Toggle(t.key)
}

func (t *Toggle) View() string {
	if global.Toggles[t.key] {
		return t.onView
	}
	return t.offView
}
