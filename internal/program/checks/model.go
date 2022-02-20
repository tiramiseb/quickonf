package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

type Model struct {
	groups      []*instructions.Group
	cursorPos   int
	viewportPos int

	width        int
	height       int
	lineToGroup  []int // when groups are filtered, a line index is not the same as the group index
	completeView []string
}

func New(groups []*instructions.Group) *Model {
	return &Model{
		groups: groups,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			return m, m.up
		case "down":
			return m, m.down
		case "pgup":
			return m, m.pgup
		case "pgdown":
			return m, m.pgdown
		case "home":
			return m, m.home
		case "end":
			return m, m.end
		}
	}
	return m, cmd
}

func (m *Model) Resize(size tea.WindowSizeMsg) (*Model, tea.Cmd) {
	m.width = size.Width
	m.height = size.Height
	return m.RedrawView()
}
