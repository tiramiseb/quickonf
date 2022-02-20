package details

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

type Model struct {
	groups []*instructions.Group

	style lipgloss.Style

	displayedGroup int

	completeView string
}

func New(groups []*instructions.Group) *Model {
	return &Model{
		groups: groups,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	// switch msg := msg.(type) {
	// case tea.KeyMsg:
	// }
	return m, cmd
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.style = lipgloss.NewStyle().Width(size.Width).Height(size.Height)
	return m
}
