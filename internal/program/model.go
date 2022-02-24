package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/checks"
	"github.com/tiramiseb/quickonf/internal/program/details"
	"github.com/tiramiseb/quickonf/internal/program/global"
	"github.com/tiramiseb/quickonf/internal/program/titlebar"
)

type model struct {
	titlebar *titlebar.Model
	checks   *checks.Model
	details  *details.Model

	leftTitle           string
	leftTitleWithFocus  string
	rightTitle          string
	rightTitleWithFocus string
	verticalSeparator   string
	subtitlesSeparator  string

	largestGroupName int

	byPriority             [][]int
	nextPriorityGroup      int
	currentlyRunningChecks int

	signalTarget chan bool
}

func newModel() *model {
	var largestName int
	for _, g := range global.AllGroups {
		l := len(g.Name)
		if l > largestName {
			largestName = l
		}
	}
	return &model{
		titlebar: titlebar.New(),
		checks:   checks.New(),
		details:  details.New(),

		largestGroupName: largestName,

		byPriority: checksIndexByPriority(),

		signalTarget: make(chan bool),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	return tea.Batch(
		m.listenSignal,
		m.next(),
	)
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
		if global.Toggles["help"] {
			switch msg.String() {
			case "ctrl+c":
				cmd = tea.Quit
			case "esc":
				global.Toggles.Disable("help")
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc", "q", "Q":
				cmd = tea.Quit
			case "f", "F":
				global.Toggles.Toggle("filter")
			case "d", "D":
				global.Toggles.Toggle("details")
			case "h", "H":
				global.Toggles.Enable("help")
			case "right":
				global.Toggles.Enable("focusOnDetails")
			case "left":
				global.Toggles.Disable("focusOnDetails")
			case "enter":
				cmd = apply(global.SelectedGroup)
			default:
				if global.Toggles["focusOnDetails"] {
					m.details, cmd = m.details.Update(msg)
				} else {
					m.checks, cmd = m.checks.Update(msg)
				}
			}
		}
	case checkDone:
		global.GroupsMayHaveChanged()
		m.currentlyRunningChecks--
		if m.currentlyRunningChecks == 0 {
			if cmd == nil {
				cmd = m.next()
			} else {
				cmd = tea.Batch(cmd, m.next())
			}
		}
	case newSignal:
		cmd = m.listenSignal
	}

	return m, cmd
}

func (m *model) View() string {
	if global.Toggles["help"] {
		return m.titlebar.View() + "\n" + m.helpView()
	} else {
		return m.titlebar.View() + "\n" + m.view()
	}
}

func (m *model) view() string {
	var leftTitle, rightTitle string
	if global.Toggles["focusOnDetails"] {
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
