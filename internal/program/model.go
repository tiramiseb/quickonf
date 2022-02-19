package program

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/checks"
	"github.com/tiramiseb/quickonf/internal/program/details"
	"github.com/tiramiseb/quickonf/internal/program/global"
	"github.com/tiramiseb/quickonf/internal/program/titlebar"
)

type model struct {
	titlebar *titlebar.Model
	checks   *checks.Model
	details  *details.Model

	groups []*instructions.Group

	byPriority             [][]int
	nextPriorityGroup      int
	currentlyRunningChecks int
}

func newModel(g []*instructions.Group) *model {
	return &model{
		titlebar: titlebar.New(),
		checks:   checks.New(g),
		details:  details.New(),

		groups: g,

		byPriority: orderChecksByPriority(g),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	return m.next()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.titlebar = m.titlebar.Resize(msg)
		width := msg.Width - 1
		height := msg.Height - 3
		left := tea.WindowSizeMsg{
			Width:  width / 2,
			Height: height,
		}
		right := tea.WindowSizeMsg{
			Width:  width - left.Width,
			Height: height,
		}
		m.checks.Resize(left)
		m.details.Resize(right)
	case tea.MouseMsg:
		if msg.Type == tea.MouseRelease {
			switch msg.Y {
			case 0:
				m.titlebar, cmd = m.titlebar.Update(msg)
			}
		}
	case tea.KeyMsg:
		if global.Global.Get("help") {
			switch msg.String() {
			case "ctrl+c":
				cmd = tea.Quit
			case "esc":
				global.Global.Set("help", false)
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc", "q", "Q":
				cmd = tea.Quit
			case "h", "H":
				global.Global.Set("help", true)
			case "f", "F":
				global.Global.Set("filter", !global.Global.Get("filter"))
			}
		}
	case checkDone:
		m.currentlyRunningChecks--
		if m.currentlyRunningChecks == 0 {
			cmd = m.next()
		}
		m.checks = m.checks.RedrawGroup(msg.groupIndex)
	}
	return m, cmd
}

func (m *model) View() string {
	if global.Global.Get("help") {
		return m.titlebar.View() + "\n" + m.helpView()
	} else {
		return m.titlebar.View() + "\n" + m.view()
	}
}

func (m *model) view() string {
	left := m.checks.View()
	right := m.details.View()
	leftWidth := lipgloss.Width(left)
	rightWidth := lipgloss.Width(right)
	leftTitle := subtitleStyle.Width(leftWidth).Render("Checks")
	rightTitle := subtitleStyle.Width(rightWidth).Render("Details")
	return leftTitle + "│" + rightTitle + "\n" +
		strings.Repeat("─", leftWidth) + "┼" + strings.Repeat("─", rightWidth) + "\n" +
		lipgloss.JoinHorizontal(lipgloss.Top, m.checks.View(), "O", m.details.View())
}

func (m *model) helpView() string {
	// TODO
	return "THERE WILL BE HELP"
}
