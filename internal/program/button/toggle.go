package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

type Toggle struct {
	actionOn  tea.Cmd
	actionOff tea.Cmd
	key       string

	offView string
	onView  string
	Width   int
}

func NewToggle(text string, hintPosition int, actionOn tea.Cmd, actionOff tea.Cmd, key string) *Toggle {
	return &Toggle{
		actionOn:  actionOn,
		actionOff: actionOff,
		key:       key,

		offView: lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		onView:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, activeHintStyle, activeStyle),
		Width:   len(text) + 2,
	}
}

func (t *Toggle) Click() (*Toggle, tea.Cmd) {
	if global.Global.Get(t.key) {
		global.Global.Set(t.key, false)
		return t, t.actionOff
	}
	global.Global.Set(t.key, true)
	return t, t.actionOn
}

func (t *Toggle) View() string {
	if global.Global.Get(t.key) {
		return t.onView
	}
	return t.offView
}
