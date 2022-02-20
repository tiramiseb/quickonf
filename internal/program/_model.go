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
	}
	return m, cmd
}
