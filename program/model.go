package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/conf"
	"github.com/tiramiseb/quickonf/instructions"
	"github.com/tiramiseb/quickonf/program/details"
	"github.com/tiramiseb/quickonf/program/groups"
	"github.com/tiramiseb/quickonf/program/help"
	"github.com/tiramiseb/quickonf/program/messages"
	"github.com/tiramiseb/quickonf/program/reallyapplyall"
	"github.com/tiramiseb/quickonf/program/reallyreloadconfig"
	"github.com/tiramiseb/quickonf/program/titlebar"
)

type model struct {
	config string
	groups *instructions.Groups

	titlebar           *titlebar.Model
	groupsview         *groups.Model
	details            *details.Model
	reallyApplyAll     *reallyapplyall.Model
	reallyReloadConfig *reallyreloadconfig.Model
	help               *help.Model

	leftTitle                       string
	leftTitleWithFocus              string
	rightTitle                      string
	rightTitleWithFocus             string
	reallyApplyRightTitle           string
	reallyApplyRightTitleWithFocus  string
	reallyReloadRightTitle          string
	reallyReloadRightTitleWithFocus string
	verticalSeparator               string
	subtitlesSeparator              string
	separatorXPos                   int

	askIfReallyApplyAll     bool
	askIfReallyReloadConfig bool
	isHelpDisplayed         bool
	focusOnDetails          bool

	size tea.WindowSizeMsg

	signalTarget chan bool
}

func newModel(config string) *model {
	d := details.New()
	return &model{
		config:             config,
		titlebar:           titlebar.New(),
		groupsview:         groups.New(d),
		details:            d,
		reallyApplyAll:     reallyapplyall.New(),
		reallyReloadConfig: reallyreloadconfig.New(),
		help:               help.New(),

		signalTarget: make(chan bool),
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.listenSignal,
		m.doReloadConfig,
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
		case m.isHelpDisplayed:
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
			// Click on detail or "really apply all?"
			msg.Y -= 3
			msg.X -= m.separatorXPos + 1
			if m.askIfReallyApplyAll {
				m.reallyApplyAll, cmd = m.reallyApplyAll.Update(msg)
				if msg.Type == tea.MouseRelease {
					m.focusOnDetails = true
				}
			} else if m.askIfReallyReloadConfig {
				m.reallyReloadConfig, cmd = m.reallyReloadConfig.Update(msg)
				if msg.Type == tea.MouseRelease {
					m.focusOnDetails = true
				}
			} else {
				m.details, cmd = m.details.Update(msg)
				if msg.Type == tea.MouseRelease {
					m.focusOnDetails = true
				}
			}
		}
	case tea.KeyMsg:
		switch {
		case m.isHelpDisplayed:
			switch msg.String() {
			case "ctrl+c":
				cmd = tea.Quit
			case "esc":
				m.isHelpDisplayed = false
			default:
				m.help, cmd = m.help.Update(msg)
			}
		case m.askIfReallyApplyAll:
			switch msg.String() {
			case "ctrl+c", "q", "Q":
				cmd = tea.Quit
			case "esc":
				cmd = m.toggleApplyAll
			case "enter", "y", "Y":
				cmd = m.doApplyAll
			case "n", "N", "l", "L":
				cmd = m.toggleApplyAll
			}
		case m.askIfReallyReloadConfig:
			switch msg.String() {
			case "ctrl+c", "q", "Q":
				cmd = tea.Quit
			case "esc":
				cmd = m.toggleReloadConfig
			case "enter", "y", "Y":
				cmd = m.doReloadConfig
			case "n", "N", "l", "L":
				cmd = m.toggleReloadConfig
			}
		default:
			switch msg.String() {
			case "ctrl+c", "esc", "q", "Q":
				cmd = tea.Quit
			case "c", "C":
				cmd = m.toggleReloadConfig
			case "f", "F":
				cmd = m.groupsview.ToggleShowSuccessful
			case "d", "D":
				cmd = m.details.ToggleDetails
			case "h", "H", "?":
				m.isHelpDisplayed = true
			case "l", "L":
				cmd = m.toggleApplyAll
			case "right":
				m.focusOnDetails = true
			case "left":
				m.focusOnDetails = false
			case "r", "R":
				cmd = m.groupsview.RecheckSelected(m.signalTarget)
			case "enter", "a", "A":
				cmd = m.groupsview.ApplySelected
			default:
				if m.focusOnDetails {
					m.details, cmd = m.details.Update(msg)
				} else {
					m.groupsview, cmd = m.groupsview.Update(msg)
				}
			}
		}
	case messages.NewSignal:
		m.groupsview, cmd = m.groupsview.Update(msg)
		cmd = tea.Batch(cmd, m.listenSignal)
	case messages.Apply:
		cmd = m.groupsview.ApplySelected
	case messages.ApplyAll:
		cmd = m.toggleApplyAll
	case messages.ConfirmApplyAll:
		cmd = m.doApplyAll
	case messages.Filter:
		cmd = m.groupsview.ToggleShowSuccessful
	case messages.Details:
		cmd = m.details.ToggleDetails
	case messages.Help:
		m.isHelpDisplayed = !m.isHelpDisplayed
	case messages.Recheck:
		cmd = m.groupsview.RecheckSelected(m.signalTarget)
	case messages.ReloadConfig:
		cmd = m.toggleReloadConfig
	case messages.ConfirmReloadConfig:
		cmd = m.doReloadConfig
	case messages.ToggleStatus:
		m.titlebar, cmd = m.titlebar.Update(msg)
	}
	return m, cmd
}

func (m *model) View() string {
	switch {
	case m.isHelpDisplayed:
		return m.titlebar.HelpView() + "\n" + m.help.View()
	default:
		return m.titlebar.View() + "\n" + m.view()
	}
}

func (m *model) view() string {
	var leftTitle, rightTitle, rightView string
	if m.focusOnDetails {
		leftTitle = m.leftTitle
		if m.askIfReallyApplyAll {
			rightTitle = m.reallyApplyRightTitleWithFocus
		} else if m.askIfReallyReloadConfig {
			rightTitle = m.reallyReloadRightTitleWithFocus
		} else {
			rightTitle = m.rightTitleWithFocus
		}
	} else {
		leftTitle = m.leftTitleWithFocus
		if m.askIfReallyApplyAll {
			rightTitle = m.reallyApplyRightTitle
		} else if m.askIfReallyReloadConfig {
			rightTitle = m.reallyReloadRightTitle
		} else {
			rightTitle = m.rightTitle
		}
	}
	if m.askIfReallyApplyAll {
		rightView = m.reallyApplyAll.View()
	} else if m.askIfReallyReloadConfig {
		rightView = m.reallyReloadConfig.View()
	} else {
		rightView = m.details.View()
	}
	return leftTitle + "│" + rightTitle + "\n" + m.subtitlesSeparator + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, m.groupsview.View(), m.verticalSeparator, rightView)
}

func (m *model) toggleApplyAll() tea.Msg {
	m.askIfReallyApplyAll = !m.askIfReallyApplyAll
	return messages.ToggleStatus{Name: "applyall", Status: m.askIfReallyApplyAll}
}

func (m *model) doApplyAll() tea.Msg {
	m.groups.ApplyAll()
	m.askIfReallyApplyAll = false
	return messages.ToggleStatus{Name: "applyall", Status: false}
}

func (m *model) toggleReloadConfig() tea.Msg {
	m.askIfReallyReloadConfig = !m.askIfReallyReloadConfig
	return messages.ToggleStatus{Name: "reloadconfig", Status: m.askIfReallyReloadConfig}
}

func (m *model) doReloadConfig() tea.Msg {
	g, errs := conf.ReadFile(m.config)
	if errs != nil {
		return tea.Quit
	}
	m.groups = g
	m.groupsview.ReplaceGroups(g)
	m.refreshSize()
	m.askIfReallyReloadConfig = false
	go m.groups.InitialChecks(m.signalTarget)
	return messages.ToggleStatus{Name: "reloadconfig", Status: false}
}
