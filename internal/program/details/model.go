package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	style  lipgloss.Style
	width  int
	height int

	completeView string
}

func New() *Model {
	return &Model{}
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.style = lipgloss.NewStyle().Width(m.width)
	m.width = size.Width
	m.height = size.Height
	m.prepareView()
	return m
}
