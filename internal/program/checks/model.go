package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

type Model struct {
	style lipgloss.Style
	// width  int
	// height int

	groups []*instructions.Group

	completeView []string
}

func New(groups []*instructions.Group) *Model {
	return &Model{
		groups: groups,
	}
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.style = lipgloss.NewStyle().Width(size.Width).Height(size.Height)
	// m.width = size.Width
	// m.height = size.Height
	m.prepareView(size.Width)
	return m
}
