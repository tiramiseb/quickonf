package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/details"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
	"github.com/tiramiseb/quickonf/internal/program/groups"
	"github.com/tiramiseb/quickonf/internal/program/help"
	"github.com/tiramiseb/quickonf/internal/program/messages"
	"github.com/tiramiseb/quickonf/internal/program/titlebar"
)

type model struct {
	groups *instructions.Groups

	titlebar   *titlebar.Model
	groupsview *groups.Model
	details    *details.Model
	help       *help.Model

	leftTitle           string
	leftTitleWithFocus  string
	rightTitle          string
	rightTitleWithFocus string
	verticalSeparator   string
	subtitlesSeparator  string
	separatorXPos       int

	focusOnDetails bool

	signalTarget chan bool
}

func newModel(g *instructions.Groups) *model {
	initialGroup := g.FirstGroup()
	details := details.New(initialGroup)
	return &model{
		groups:     g,
		titlebar:   titlebar.New(),
		groupsview: groups.New(g, initialGroup, details),
		details:    details,
		help:       help.New(),

		signalTarget: make(chan bool),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	return tea.Batch(
		m.listenSignal,
		func() tea.Msg {
			m.groups.InitialChecks(m.signalTarget)
			return nil
		},
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
			m.groupsview, cmd = m.groupsview.Update(msg)
			if msg.Type == tea.MouseRelease {
				m.focusOnDetails = false
			}
		case msg.X > m.separatorXPos:
			// Clock on detail
			msg.Y -= 3
			m.details, cmd = m.details.Update(msg)
			if msg.Type == tea.MouseRelease {
				m.focusOnDetails = true
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
				m.focusOnDetails = true
			case "left":
				m.focusOnDetails = false
			// case "r", "R":
			//
			// case "enter":
			// 	go groups.ApplySelected()
			default:
				if m.focusOnDetails {
					m.details, cmd = m.details.Update(msg)
				} else {
					m.groupsview, cmd = m.groupsview.Update(msg)
				}
			}
		}
	case messages.SelectedGroup:
		m.details, cmd = m.details.Update(msg)
	case messages.NewSignal:
		m.groupsview, cmd = m.groupsview.Update(msg)
		cmd = tea.Batch(cmd, m.listenSignal)
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
	if m.focusOnDetails {
		leftTitle = m.leftTitle
		rightTitle = m.rightTitleWithFocus
	} else {
		leftTitle = m.leftTitleWithFocus
		rightTitle = m.rightTitle
	}
	return leftTitle + "│" + rightTitle + "\n" + m.subtitlesSeparator + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, m.groupsview.View(), m.verticalSeparator, m.details.View())
}
