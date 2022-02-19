package program

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/checks"
	"github.com/tiramiseb/quickonf/internal/program/details"
	"github.com/tiramiseb/quickonf/internal/program/titlebar"
)

type model struct {
	titlebar *titlebar.Model
	checks   *checks.Model
	details  *details.Model

	showHelp bool
}

func newModel(g []*instructions.Group) *model {
	return &model{
		titlebar: titlebar.New(),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.titlebar = m.titlebar.Resize(msg)
		// TODO Calculate left width and right width, for check and details
		// No need to store it here, just pass it to check and details
	case tea.MouseMsg:
		if msg.Type == tea.MouseRelease {
			switch msg.Y {
			case 0:
				m.titlebar, cmd = m.titlebar.Update(msg)
			}
		}
	case tea.KeyMsg:
		if m.showHelp {
			switch msg.String() {
			case "ctrl+c":
				cmd = tea.Quit
			case "esc":
				m.showHelp = false
				m.titlebar = m.titlebar.HelpActive(false)
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc", "q", "Q":
				cmd = tea.Quit
			case "h", "H":
				m.showHelp = true
				m.titlebar = m.titlebar.HelpActive(true)
			}
		}
	}
	return m, cmd
}

func (m *model) View() string {
	var content string
	if m.showHelp {
		content = m.helpView()
	} else {
		content = m.view()
	}
	return m.titlebar.View() + "\n" + content
}

func (m *model) view() string {
	left := m.checks.View()
	right := m.details.View()
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	leftTitle := subtitleStyle.Width(leftWidth).Render("Checks")
	rightTitle := subtitleStyle.Width(rightWidth).Render("Details")
	return leftTitle + "│" + rightTitle + "\n" +
		strings.Repeat("─", leftWidth) + "┼" + strings.Repeat("─", rightWidth)

	// content = lipgloss.JoinHorizontal(
	// 	lipgloss.Top,
	// 	m.checks.View(),
	// 	// m.separator.View(),
	// 	m.details.View(),
	// )
}

func (m *model) helpView() string {
	return ""
}
