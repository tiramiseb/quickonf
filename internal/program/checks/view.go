package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) RedrawView() (*Model, tea.Cmd) {
	if m.width == 0 {
		return m, nil
	}
	filter := global.Toggles["filter"].Get()
	var (
		view        []string
		lineToGroup []int
	)
	var selectedGroup int
	if m.cursorPos < len(m.lineToGroup) {
		selectedGroup = m.lineToGroup[m.cursorPos]
	}
	for i, g := range m.groups {
		status := g.Status()
		if filter && status == commands.StatusSuccess {
			// If selected group is hidden, select the next group
			if i == selectedGroup {
				selectedGroup++
			}
			continue
		}
		if i == selectedGroup {
			m.cursorPos = len(view)
		}
		view = append(view, global.MakeWidth(g.Name, m.width))
		lineToGroup = append(lineToGroup, i)
	}

	m.completeView = view
	m.lineToGroup = lineToGroup
	m.checkCursorVsViewport()
	return m, func() tea.Msg {
		return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
	}
}

func (m *Model) checkCursorVsViewport() {
	if m.cursorPos < 0 {
		m.cursorPos = 0
	}
	if m.cursorPos >= len(m.completeView) {
		m.cursorPos = len(m.completeView) - 1
	}
	if m.cursorPos < m.viewportPos {
		m.viewportPos = m.cursorPos
	}
	if m.cursorPos >= m.viewportPos+m.height {
		m.viewportPos = m.cursorPos - m.height + 1
	}
}

func (m *Model) View() string {
	if len(m.completeView) == 0 {
		return "No check"
	}
	var view string
	lastLine := m.viewportPos + m.height - 1
	if lastLine >= len(m.completeView)-1 {
		lastLine = len(m.completeView) - 1
	}
	for i := m.viewportPos; i < lastLine; i++ {
		if i == m.cursorPos {
			view += global.SelectedStyles[m.groups[m.lineToGroup[i]].Status()].Render(m.completeView[i]) + "\n"
		} else {
			view += global.Styles[m.groups[m.lineToGroup[i]].Status()].Render(m.completeView[i]) + "\n"
		}
	}
	if lastLine == m.cursorPos {
		view += global.SelectedStyles[m.groups[m.lineToGroup[lastLine]].Status()].Render(m.completeView[lastLine])
	} else {
		view += global.Styles[m.groups[m.lineToGroup[lastLine]].Status()].Render(m.completeView[lastLine])
	}
	return view
}
