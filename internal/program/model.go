package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/specific/applys"
	"github.com/tiramiseb/quickonf/internal/program/specific/checks"
	"github.com/tiramiseb/quickonf/internal/program/specific/separator"
	"github.com/tiramiseb/quickonf/internal/program/specific/titlebar"
)

type model struct {
	titlebar  tea.Model
	checks    tea.Model
	separator tea.Model
	applys    tea.Model

	leftPartEndColumn    int
	rightPartStartColumn int
	activeApply          bool // If false, the "check" part is active
}

func newModel(g []*instructions.Group) *model {
	return &model{
		titlebar:  titlebar.New(),
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
		cmds := make([]tea.Cmd, 4)
		m.titlebar, cmds[0] = m.titlebar.Update(
			tea.WindowSizeMsg{Height: 1, Width: msg.Width},
		)
		leftWidth := (msg.Width - 1) / 2
		rightWidth := msg.Width - leftWidth - 1
		m.checks, cmds[1] = m.checks.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: leftWidth},
		)
		m.applys, cmds[2] = m.applys.Update(
			tea.WindowSizeMsg{Height: msg.Height - 1, Width: rightWidth},
		)
		m.leftPartEndColumn = leftWidth - 1
		m.rightPartStartColumn = leftWidth + 1
		m.separator, cmds[3] = m.separator.Update(tea.WindowSizeMsg{Height: msg.Height - 1, Width: 1})
		cmd = tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			cmd = m.activateApplies()
		case "left":
			cmd = m.activateChecks()
		default:
			m.titlebar, cmd = m.titlebar.Update(msg)
			if cmd == nil {
				if m.activeApply {
					m.applys, cmd = m.applys.Update(msg)
				} else {
					m.checks, cmd = m.checks.Update(msg)
				}
			}
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
		cmd = tea.Batch(cmds...)
	case separator.CursorMsg:
		m.separator, cmd = m.separator.Update(msg)
	default:
		cmds := make([]tea.Cmd, 3)
		m.titlebar, cmds[0] = m.titlebar.Update(msg)
		m.checks, cmds[1] = m.checks.Update(msg)
		m.applys, cmds[2] = m.applys.Update(msg)
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
	return m.titlebar.View() + "\n" +
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.checks.View(),
			m.separator.View(),
			m.applys.View(),
		)
}
