package groups

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	selectedGroupToViewportOffset int // How much must be removed from selected group to get viewport start

	width  int
	height int
}

func New() *Model {
	return &Model{}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.up()
		case "down":
			m.down()
		case "pgup":
			m.pgup()
		case "pgdown":
			m.pgdown()
		case "home":
			m.home()
		case "end":
			m.end()
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.scrollUp()
		case tea.MouseWheelDown:
			m.scrollDown()
		case tea.MouseRelease:
			m.selectLine(msg.Y)
		}
	}
	return m, cmd
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.width = size.Width
	m.height = size.Height
	return m
}
