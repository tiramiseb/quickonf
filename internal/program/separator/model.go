package separator

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

type CursorMsg struct {
	PointingApply bool
	Position      int
}

type model struct {
	height        int
	currentCursor CursorMsg
}

func New() *model {
	return &model{height: 4}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
	case CursorMsg:
		m.currentCursor = msg
	}
	return m, nil
}

func (m *model) View() string {
	var cursor string
	if m.currentCursor.PointingApply {
		cursor = "⏵"
	} else {
		cursor = "⏴"
	}
	return style.BoxContent.Copy().
		Height(m.height).
		Bold(true).
		Render(
			strings.Repeat("\n", m.currentCursor.Position+2) + cursor,
		)
}
