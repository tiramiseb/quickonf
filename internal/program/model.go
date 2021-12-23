package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/groupapplys"
	"github.com/tiramiseb/quickonf/internal/program/groupchecks"
	"github.com/tiramiseb/quickonf/internal/program/separator"
	"github.com/tiramiseb/quickonf/internal/program/titlebar"
)

type model struct {
	titlebar  tea.Model
	checks    tea.Model
	separator tea.Model
	applys    tea.Model

	leftPartEndColumn    int
	rightPartStartColumn int
	activeApply          bool // If false, "check" is active
}

func newModel(g []*instructions.Group) *model {
	return &model{
		titlebar:  titlebar.New(),
		checks:    groupchecks.New(g),
		applys:    groupapplys.New(),
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
			m.checks, _ = m.checks.Update(groupchecks.ActiveMsg{Active: false})
			m.applys, _ = m.applys.Update(groupapplys.ActiveMsg{Active: true})
			m.activeApply = true
		case "left":
			m.checks, _ = m.checks.Update(groupchecks.ActiveMsg{Active: true})
			m.applys, _ = m.applys.Update(groupapplys.ActiveMsg{Active: false})
			m.activeApply = false
		}
		m.titlebar, cmd = m.titlebar.Update(msg)
		if cmd == nil {
			if m.activeApply {
				m.applys, cmd = m.applys.Update(msg)
			} else {
				m.checks, cmd = m.checks.Update(msg)
			}
		}
	case tea.MouseMsg:
		unknown := tea.MouseMsg{
			Type: tea.MouseUnknown,
		}
		switch msg.Y {
		case 0:
			m.titlebar, cmd = m.titlebar.Update(msg)
			m.applys, _ = m.applys.Update(unknown)
			m.applys, _ = m.applys.Update(unknown)
		default:
			m.titlebar, _ = m.titlebar.Update(unknown)
			msg.Y--
			if msg.X <= m.leftPartEndColumn {
				m.checks, cmd = m.checks.Update(msg)
				m.applys, _ = m.applys.Update(unknown)
			} else if msg.X >= m.leftPartEndColumn {
				msg.X -= m.rightPartStartColumn
				m.applys, cmd = m.applys.Update(msg)
				m.checks, _ = m.checks.Update(unknown)
			}
		}
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

func (m *model) View() string {
	return m.titlebar.View() + "\n" +
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.checks.View(),
			m.separator.View(),
			m.applys.View(),
		)
}
