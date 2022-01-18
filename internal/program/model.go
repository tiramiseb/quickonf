package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/messages"
	"github.com/tiramiseb/quickonf/internal/program/specific/applys"
	"github.com/tiramiseb/quickonf/internal/program/specific/checks"
	"github.com/tiramiseb/quickonf/internal/program/specific/help"
	"github.com/tiramiseb/quickonf/internal/program/specific/separator"
	"github.com/tiramiseb/quickonf/internal/program/specific/titlebar"
)

type model struct {
	titlebar  tea.Model
	help      tea.Model
	checks    tea.Model
	separator tea.Model
	applys    tea.Model

	leftPartEndColumn    int
	rightPartStartColumn int
	activeApply          bool // If false, the "check" part is active
	showHelp             bool
	filtered             bool
}

func newModel(g []*instructions.Group) *model {
	return &model{
		titlebar:  titlebar.New(),
		help:      help.New(),
		checks:    checks.New(g),
		applys:    applys.New(),
		separator: separator.New(),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	return m.checks.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmds := make([]tea.Cmd, 5)
		m.titlebar, cmds[0] = m.titlebar.Update(
			tea.WindowSizeMsg{Height: 1, Width: msg.Width},
		)
		m.help, cmds[1] = m.help.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: msg.Width},
		)
		leftWidth := (msg.Width - 1) / 2
		rightWidth := msg.Width - leftWidth - 1
		m.checks, cmds[2] = m.checks.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: leftWidth},
		)
		m.applys, cmds[3] = m.applys.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: rightWidth},
		)
		m.leftPartEndColumn = leftWidth - 1
		m.rightPartStartColumn = leftWidth + 1
		m.separator, cmds[4] = m.separator.Update(tea.WindowSizeMsg{Height: msg.Height - 1, Width: 1})
		cmd = tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			cmd = tea.Quit
		case "esc":
			if m.showHelp {
				m.showHelp = false
				cmd1 := messages.Help(m.showHelp)
				var cmd2 tea.Cmd
				m.titlebar, cmd2 = m.titlebar.Update(cmd1())
				cmd = tea.Batch(cmd1, cmd2)
			} else {
				cmd = tea.Quit
			}
		case "q", "Q":
			if m.showHelp {
				m.help, cmd = m.help.Update(msg)
			} else {
				cmd = tea.Quit
			}
		case "h", "H":
			if m.showHelp {
				m.help, cmd = m.help.Update(msg)
			} else {
				m.showHelp = !m.showHelp
				cmd1 := messages.Help(m.showHelp)
				var cmd2 tea.Cmd
				m.titlebar, cmd2 = m.titlebar.Update(cmd1())
				cmd = tea.Batch(cmd1, cmd2)
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
			case m.activeApply:
				m.applys, cmd = m.applys.Update(msg)
			default:
				m.checks, cmd = m.checks.Update(msg)
			}
		}
	case messages.ToggleMsg:
		if m.activeApply {
			m.applys, cmd = m.applys.Update(msg)
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
		case 0:
			m.titlebar, subCmd = m.titlebar.Update(msg)
			if subCmd != nil {
				cmds = append(cmds, subCmd)
			}
			m.checks, subCmd = m.checks.Update(unknown)
			if subCmd != nil {
				cmds = append(cmds, subCmd)
			}
			m.applys, subCmd = m.applys.Update(unknown)
			if subCmd != nil {
				cmds = append(cmds, subCmd)
			}
		default:
			m.titlebar, subCmd = m.titlebar.Update(unknown)
			if subCmd != nil {
				cmds = append(cmds, subCmd)
			}
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
					m.applys, subCmd = m.applys.Update(unknown)
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
					m.applys, subCmd = m.applys.Update(msg)
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
					m.applys, subCmd = m.applys.Update(unknown)
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
		m.applys, cmds[3] = m.applys.Update(msg)
		cmd = tea.Batch(cmds...)
	}
	return m, cmd
}

func (m *model) activateApplies() tea.Cmd {
	var cmd1, cmd2, cmd3 tea.Cmd
	activeMsg := separator.ActiveMsg{IsRightActive: true}
	m.titlebar, cmd1 = m.titlebar.Update(activeMsg)
	m.checks, cmd2 = m.checks.Update(activeMsg)
	m.applys, cmd3 = m.applys.Update(activeMsg)
	m.activeApply = true
	return tea.Batch(cmd1, cmd2, cmd3)
}

func (m *model) activateChecks() tea.Cmd {
	var cmd1, cmd2, cmd3 tea.Cmd
	activeMsg := separator.ActiveMsg{IsRightActive: false}
	m.titlebar, cmd1 = m.titlebar.Update(activeMsg)
	m.checks, cmd2 = m.checks.Update(activeMsg)
	m.applys, cmd3 = m.applys.Update(activeMsg)
	m.activeApply = false
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
			m.applys.View(),
		)
	}
	return m.titlebar.View() + "\n" + content
}
