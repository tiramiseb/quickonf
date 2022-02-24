package details

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	viewport viewport.Model

	width int
}

func New() *Model {
	return &Model{
		viewport: viewport.Model{Width: 1, Height: 1},
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.viewport.Height = size.Height
	m.viewport.Width = size.Width
	m.width = size.Width
	return m
}
