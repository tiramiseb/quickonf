package details

import tea "github.com/charmbracelet/bubbletea"

type Model struct{}

func New() *Model {
	return &Model{}
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.draw(size.Width, size.Height)
	return m
}
