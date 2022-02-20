package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

type Toggle struct {
	trigger tea.Cmd
	key     string

	offView string
	onView  string
	Width   int
}

func NewToggle(text string, hintPosition int, trigger tea.Cmd, key string) *Toggle {
	return &Toggle{
		trigger: trigger,
		key:     key,

		offView: lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		onView:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, activeHintStyle, activeStyle),
		Width:   len(text) + 2,
	}
}

func (t *Toggle) Click() (*Toggle, tea.Cmd) {
	global.Global.Set(t.key, !global.Global.Get(t.key))
	return t, t.trigger
}

func (t *Toggle) View() string {
	if global.Global.Get(t.key) {
		return t.onView
	}
	return t.offView
}
