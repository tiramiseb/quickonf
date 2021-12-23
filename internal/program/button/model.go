package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

type model struct {
	view        string
	clickedView string

	clicked bool
	action  tea.Cmd
}

func New(text string, hotkeyPos int, action tea.Cmd) *model {
	return &model{
		view:        lipgloss.StyleRunes("["+text+"]", []int{hotkeyPos + 1}, style.ButtonKey, style.Button),
		clickedView: lipgloss.StyleRunes("["+text+"]", []int{hotkeyPos + 1}, style.ClickedButtonKey, style.ClickedButton),
		action:      action,
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
				return m, m.action
			}
		}

	}
	return m, nil
}

func (m *model) View() string {
	if m.clicked {
		return m.clickedView
	}
	return m.view
}
