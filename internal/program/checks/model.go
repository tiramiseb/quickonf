package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

type Model struct {
	groups []*instructions.Group

	width        int
	height       int
	completeView string
}

func New(groups []*instructions.Group) *Model {
	return &Model{
		groups: groups,
	}
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.width = size.Width
	m.height = size.Height
	return m.RedrawView()
}
