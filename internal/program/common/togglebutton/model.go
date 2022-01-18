package togglebutton

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

type Toggle struct {
	On bool
}

type model struct {
	offView string
	onView  string

	clicked   bool
	on        bool
	actionOn  tea.Cmd
	actionOff tea.Cmd
}

func New(text string, hotkeyPos int, actionOn tea.Cmd, actionOff tea.Cmd) *model {
	return &model{
		offView:   lipgloss.StyleRunes("["+text+"]", []int{hotkeyPos + 1}, style.ButtonKey, style.Button),
		onView:    lipgloss.StyleRunes("["+text+"]", []int{hotkeyPos + 1}, style.ClickedButtonKey, style.ClickedButton),
		actionOn:  actionOn,
		actionOff: actionOff,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			m.clicked = true
		case tea.MouseUnknown:
			m.clicked = false
		case tea.MouseRelease:
			if m.clicked {
				m.clicked = false
				if m.on {
					m.on = false
					return m, m.actionOff
				}
				m.on = true
				return m, m.actionOn
			}
		}
	case Toggle:
		m.on = msg.On
	}
	return m, nil
}

func (m *model) View() string {
	if m.clicked || m.on {
		return m.onView
	}
	return m.offView
}
