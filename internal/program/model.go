package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/details"
	"github.com/tiramiseb/quickonf/internal/program/global/groups"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
	groupsList "github.com/tiramiseb/quickonf/internal/program/groups"
	"github.com/tiramiseb/quickonf/internal/program/help"
	"github.com/tiramiseb/quickonf/internal/program/titlebar"
)

type model struct {
	titlebar *titlebar.Model
	groups   *groupsList.Model
	details  *details.Model
	help     *help.Model

	leftTitle           string
	leftTitleWithFocus  string
	rightTitle          string
	rightTitleWithFocus string
	verticalSeparator   string
	subtitlesSeparator  string

	largestGroupName int
	separatorXPos    int

	byPriority             [][]int
	nextPriorityGroup      int
	currentlyRunningChecks int

	signalTarget chan bool
}

func newModel() *model {
	return &model{
		titlebar: titlebar.New(),
		groups:   groupsList.New(),
		details:  details.New(),
		help:     help.New(),

		largestGroupName: groups.GetMaxNameLength(),

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
		switch {
		case msg.Y == 0:
			// Click on titlebar
			m.titlebar, cmd = m.titlebar.Update(msg)
		case toggles.Get("help"):
			// Help is displayed, forward to help
			msg.Y--
			m.help, cmd = m.help.Update(msg)
		case msg.X < m.separatorXPos:
			// Click on group
			msg.Y -= 3
			m.groups, cmd = m.groups.Update(msg)
			if msg.Type == tea.MouseRelease {
				toggles.Disable("focusOnDetails")
			}
		case msg.X > m.separatorXPos:
			// Clock on detail
			msg.Y -= 3
			m.details, cmd = m.details.Update(msg)
			if msg.Type == tea.MouseRelease {
				toggles.Enable("focusOnDetails")
			}
		}
	case tea.KeyMsg:
		if toggles.Get("help") {
			switch msg.String() {
			case "ctrl+c":
				cmd = tea.Quit
			case "esc":
				toggles.Disable("help")
			default:
				m.help, cmd = m.help.Update(msg)
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "esc", "q", "Q":
				cmd = tea.Quit
			case "f", "F":
				toggles.Toggle("filter")
			case "d", "D":
				toggles.Toggle("details")
			case "h", "H":
				toggles.Enable("help")
			case "right":
				toggles.Enable("focusOnDetails")
			case "left":
				toggles.Disable("focusOnDetails")
			// case "r", "R":
			//
			case "enter":
				go groups.ApplySelected()
			default:
				if toggles.Get("focusOnDetails") {
					m.details, cmd = m.details.Update(msg)
				} else {
					m.groups, cmd = m.groups.Update(msg)
				}
			}
		}
	case checkDone:
		m.currentlyRunningChecks--
		if m.currentlyRunningChecks == 0 {
			cmd = m.next()
		}
	case newSignal:
		cmd = m.listenSignal
	}

	return m, cmd
}

func (m *model) View() string {
	if toggles.Get("help") {
		return m.titlebar.View() + "\n" + m.help.View()
	} else {
		return m.titlebar.View() + "\n" + m.view()
	}
}

func (m *model) view() string {
	var leftTitle, rightTitle string
	if toggles.Get("focusOnDetails") {
		leftTitle = m.leftTitle
		rightTitle = m.rightTitleWithFocus
	} else {
		leftTitle = m.leftTitleWithFocus
		rightTitle = m.rightTitle
	}
	return leftTitle + "│" + rightTitle + "\n" + m.subtitlesSeparator + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, m.groups.View(), m.verticalSeparator, m.details.View())
}
