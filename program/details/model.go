package details

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/instructions"
	"github.com/tiramiseb/quickonf/program/messages"
)

type Model struct {
	viewport viewport.Model

	showDetails bool

	group *instructions.Group
	width int
}

func New() *Model {
	return &Model{viewport: viewport.Model{Width: 1, Height: 1}}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		default:
			m.viewport, cmd = m.viewport.Update(msg)
		}
	case tea.MouseMsg:
		if msg.Y >= 0 {
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}
	return m, cmd
}

func (m *Model) ShowGroup(g *instructions.Group) {
	m.group = g
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.viewport.Height = size.Height
	m.viewport.Width = size.Width
	m.width = size.Width
	return m
}

func (m *Model) ToggleDetails() tea.Msg {
	m.showDetails = !m.showDetails
	return messages.ToggleStatus{Name: "details", Status: m.showDetails}
}
