package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	// "github.com/tiramiseb/quickonf/internal/program/common/messages"
	// "github.com/tiramiseb/quickonf/internal/program/specific/checks"
	// "github.com/tiramiseb/quickonf/internal/program/specific/details"
	// "github.com/tiramiseb/quickonf/internal/program/specific/help"
	// "github.com/tiramiseb/quickonf/internal/program/specific/separator"
)

type model struct {
	help      tea.Model
	checks    tea.Model
	separator tea.Model
	details   tea.Model

	leftPartEndColumn    int
	rightPartStartColumn int
	activeDetails        bool // If false, the "check" part is active
	showHelp             bool
	filtered             bool
}

// newModel creates a new root model for the application
//
// groups are expected to be sorted by priority
func newModel(g []*instructions.Group) *model {
	return &model{
		// help:      help.New(),
		// checks:    checks.New(g),
		// details:   details.New(),
		// separator: separator.New(),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	m.filtered = true
	return tea.Batch(
		// messages.Filter(m.filtered),
		m.checks.Init(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmds := make([]tea.Cmd, 5)
		m.help, cmds[1] = m.help.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: msg.Width},
		)
		leftWidth := (msg.Width - 1) / 2
		rightWidth := msg.Width - leftWidth - 1
		m.checks, cmds[2] = m.checks.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: leftWidth},
		)
		m.details, cmds[3] = m.details.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: rightWidth},
		)
		m.leftPartEndColumn = leftWidth - 1
		m.rightPartStartColumn = leftWidth + 1
		m.separator, cmds[4] = m.separator.Update(tea.WindowSizeMsg{Height: msg.Height - 1, Width: 1})
		cmd = tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "h", "H":
			if m.showHelp {
				m.help, cmd = m.help.Update(msg)
			} else {
				m.showHelp = true
				cmd1 := messages.Help(m.showHelp)
				var cmd2, cmd3 tea.Cmd
				m.titlebar, cmd2 = m.titlebar.Update(cmd1())
				m.help, cmd3 = m.help.Update(cmd1())
				cmd = tea.Batch(cmd1, cmd2, cmd3)
			}
		case "f", "F":
			if m.showHelp {
				m.help, cmd = m.help.Update(msg)
			} else {
				m.filtered = !m.filtered
				cmd = messages.Filter(m.filtered)
			}
		case "right":
			if m.showHelp {
				m.help, cmd = m.help.Update(msg)
			} else {
				cmd = m.activateApplies()
			}
		case "left":
			if m.showHelp {
				m.help, cmd = m.help.Update(msg)
			} else {
				cmd = m.activateChecks()
			}
		default:
			switch {
			case m.showHelp:
				m.help, cmd = m.help.Update(msg)
			case m.activeDetails:
				m.details, cmd = m.details.Update(msg)
			default:
				m.checks, cmd = m.checks.Update(msg)
			}
		}
	case messages.ToggleMsg:
		if m.activeDetails {
			m.details, cmd = m.details.Update(msg)
		} else {
			m.checks, cmd = m.checks.Update(msg)
		}
	case tea.MouseMsg:
		unknown := tea.MouseMsg{
			Type: tea.MouseUnknown,
		}
		var subCmd tea.Cmd
		var cmds []tea.Cmd
		switch msg.Y {
		default:
			msg.Y--
			if m.showHelp {
				var subCmd tea.Cmd
				m.help, subCmd = m.help.Update(msg)
				if subCmd != nil {
					cmds = append(cmds, subCmd)
				}
			} else {
				if msg.X <= m.leftPartEndColumn {
					// Over checks
					if msg.Type != tea.MouseMotion {
						cmds = append(cmds, m.activateChecks())
					}
					m.checks, subCmd = m.checks.Update(msg)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
					m.details, subCmd = m.details.Update(unknown)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
				} else if msg.X >= m.rightPartStartColumn {
					// Over applies
					if msg.Type != tea.MouseMotion {
						cmds = append(cmds, m.activateApplies())
					}
					msg.X -= m.rightPartStartColumn
					m.checks, subCmd = m.checks.Update(unknown)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
					m.details, subCmd = m.details.Update(msg)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
				} else {
					// Over the separator
					m.titlebar, subCmd = m.titlebar.Update(unknown)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
					m.checks, subCmd = m.checks.Update(unknown)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
					m.details, subCmd = m.details.Update(unknown)
					if subCmd != nil {
						cmds = append(cmds, subCmd)
					}
				}
			}
		}
		cmd = tea.Batch(cmds...)
	case separator.CursorMsg:
		m.separator, cmd = m.separator.Update(msg)
	case messages.HelpMsg:
		m.showHelp = msg.On
	default:
		cmds := make([]tea.Cmd, 4)
		m.titlebar, cmds[0] = m.titlebar.Update(msg)
		m.help, cmds[1] = m.help.Update(msg)
		m.checks, cmds[2] = m.checks.Update(msg)
		m.details, cmds[3] = m.details.Update(msg)
		cmd = tea.Batch(cmds...)
	}
	return m, cmd
}

func (m *model) activateApplies() tea.Cmd {
	var cmd1, cmd2, cmd3 tea.Cmd
	activeMsg := separator.ActiveMsg{IsRightActive: true}
	m.checks, cmd2 = m.checks.Update(activeMsg)
	m.details, cmd3 = m.details.Update(activeMsg)
	m.activeDetails = true
	return tea.Batch(cmd1, cmd2, cmd3)
}

func (m *model) activateChecks() tea.Cmd {
	var cmd1, cmd2, cmd3 tea.Cmd
	activeMsg := separator.ActiveMsg{IsRightActive: false}
	m.checks, cmd2 = m.checks.Update(activeMsg)
	m.details, cmd3 = m.details.Update(activeMsg)
	m.activeDetails = false
	return tea.Batch(cmd1, cmd2, cmd3)
}

func (m *model) View() string {
	var content string
	if m.showHelp {
		content = m.help.View()
	} else {
		content = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.checks.View(),
			m.separator.View(),
			m.details.View(),
		)
	}
	return m.titlebar.View() + "\n" + content
}
