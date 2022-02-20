package global

import tea "github.com/charmbracelet/bubbletea"

type toggle struct {
	name  string
	value bool
}

type ToggleMsg struct {
	Name string
}

var Toggles = map[string]*toggle{
	"details":        {"detalis", false},
	"filter":         {"filter", true},
	"focusOnDetails": {"focusOnDetails", false},
	"help":           {"help", false},
}

func (t *toggle) Toggle() tea.Msg {
	t.value = !t.value
	return ToggleMsg{t.name}
}

func (t *toggle) Enable() tea.Msg {
	t.value = true
	return ToggleMsg{t.name}
}
func (t *toggle) Disable() tea.Msg {
	t.value = false
	return ToggleMsg{t.name}
}

func (t *toggle) Get() bool {
	return t.value
}
