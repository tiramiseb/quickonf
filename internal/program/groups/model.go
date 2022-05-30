package groups

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/details"
	"github.com/tiramiseb/quickonf/internal/program/messages"
)

type Model struct {
	groups  *instructions.Groups
	details *details.Model

	firstDisplayedGroup *instructions.Group
	selectedGroup       *instructions.Group

	showSuccessful bool

	width  int
	height int
}

func New(g *instructions.Groups, initialGroup *instructions.Group, d *details.Model) *Model {
	return &Model{
		groups:        g,
		details:       d,
		selectedGroup: initialGroup,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.up()
			cmd = m.selectedGroupCmd
		case "down":
			m.down()
			cmd = m.selectedGroupCmd
		case "pgup":
			m.pgup()
			cmd = m.selectedGroupCmd
		case "pgdown":
			m.pgdown()
			cmd = m.selectedGroupCmd
		case "home":
			m.home()
			cmd = m.selectedGroupCmd
		case "end":
			m.end()
			cmd = m.selectedGroupCmd
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.scrollUp()
			cmd = m.selectedGroupCmd
		case tea.MouseWheelDown:
			m.scrollDown()
			cmd = m.selectedGroupCmd
		case tea.MouseRelease:
			if msg.Y >= 0 {
				m.selectLine(msg.Y)
				cmd = m.selectedGroupCmd
			}
		}
	case messages.NewSignal:
		m.selectedGroup = m.selectedGroup.Next(0, m.showSuccessful)
		cmd = m.selectedGroupCmd
	}
	return m, cmd
}

func (m *Model) selectedGroupCmd() tea.Msg {
	return messages.SelectedGroup{Group: m.selectedGroup}
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.width = size.Width
	m.height = size.Height
	return m
}

func (m *Model) ToggleShowSuccessful() tea.Msg {
	m.showSuccessful = !m.showSuccessful
	return messages.ToggleStatus{Name: "filter", Status: !m.showSuccessful}
}
