package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Toggle struct {
	offView string
	onView  string
	Width   int

	actionOn  tea.Cmd
	actionOff tea.Cmd

	isOn bool
}

func NewToggle(text string, hintPosition int, actionOn tea.Cmd, actionOff tea.Cmd, initialValue bool) *Toggle {
	return &Toggle{
		offView: lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, inactiveHintStyle, inactiveStyle),
		onView:  lipgloss.StyleRunes("["+text+"]", []int{hintPosition + 1}, activeHintStyle, activeStyle),
		Width:   len(text) + 2,

		actionOn:  actionOn,
		actionOff: actionOff,

		isOn: initialValue,
	}
}

func (t *Toggle) Click() (*Toggle, tea.Cmd) {
	if t.isOn {
		t.isOn = false
		return t, t.actionOff
	}
	t.isOn = true
	return t, t.actionOn
}

func (t *Toggle) View() string {
	if t.isOn {
		return t.onView
	}
	return t.offView
}

func (t *Toggle) FromExternal(on bool) *Toggle {
	t.isOn = on
	return t
}
