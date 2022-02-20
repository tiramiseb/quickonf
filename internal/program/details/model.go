package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	style lipgloss.Style

	completeView string
}

func New() *Model {
	return &Model{}
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.style = lipgloss.NewStyle().Width(size.Width).Height(size.Height)
	return m.RedrawView()
}
