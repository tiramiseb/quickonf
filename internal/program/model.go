package program

import (
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

	leftTitle           string
	leftTitleWithFocus  string
	rightTitle          string
	rightTitleWithFocus string
	verticalSeparator   string
	subtitlesSeparator  string

	byPriority             [][]int
	nextPriorityGroup      int
	currentlyRunningChecks int
}

func newModel(g []*instructions.Group) *model {
	return &model{
		titlebar: titlebar.New(),
		checks:   checks.New(g),
		details:  details.New(g),

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
		m.resize(msg)
	case tea.MouseMsg:
		if msg.Type == tea.MouseRelease {
			switch msg.Y {
			case 0:
				m.titlebar, cmd = m.titlebar.Update(msg)
			}
		}
	case tea.KeyMsg:
		if global.Toggles["help"].Get() {
			switch msg.String() {
			case "ctrl+c":
				cmd = tea.Quit
			case "esc":
				cmd = global.Toggles["help"].Disable
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc", "q", "Q":
				cmd = tea.Quit
			case "f", "F":
				cmd = global.Toggles["filter"].Toggle
			case "d", "D":
				cmd = global.Toggles["details"].Toggle
			case "h", "H":
				cmd = global.Toggles["help"].Enable
			case "right":
				cmd = global.Toggles["focusOnDetails"].Enable
			case "left":
				cmd = global.Toggles["focusOnDetails"].Disable
			default:
				if global.Toggles["focusOnDetails"].Get() {
					m.details, cmd = m.details.Update(msg)
				} else {
					m.checks, cmd = m.checks.Update(msg)
				}
			}
		}
	case checkDone:
		m.checks, cmd = m.checks.RedrawView()
		m.currentlyRunningChecks--
		if m.currentlyRunningChecks == 0 {
			if cmd == nil {
				cmd = m.next()
			} else {
				cmd = tea.Batch(cmd, m.next())
			}
		}
	case global.ToggleMsg:
		switch msg.Name {
		case "filter":
			m.checks, cmd = m.checks.RedrawView()
		}
	case global.SelectGroupMsg:
		m.details = m.details.ChangeView(msg.Idx)
	}
	return m, cmd
}

func (m *model) View() string {
	if global.Toggles["help"].Get() {
		return m.titlebar.View() + "\n" + m.helpView()
	} else {
		return m.titlebar.View() + "\n" + m.view()
	}
}

func (m *model) view() string {
	var leftTitle, rightTitle string
	if global.Toggles["focusOnDetails"].Get() {
		leftTitle = m.leftTitle
		rightTitle = m.rightTitleWithFocus
	} else {
		leftTitle = m.leftTitleWithFocus
		rightTitle = m.rightTitle
	}
	return leftTitle + "│" + rightTitle + "\n" + m.subtitlesSeparator + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, m.checks.View(), m.verticalSeparator, m.details.View())
}

func (m *model) helpView() string {
	// TODO
	return "THERE WILL BE HELP"
}
