package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) up() tea.Msg {
	m.cursorPos--
	m.checkCursorVsViewport()
	return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
}

func (m *Model) down() tea.Msg {
	m.cursorPos++
	m.checkCursorVsViewport()
	return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
}

func (m *Model) pgup() tea.Msg {
	if m.cursorPos == m.viewportPos {
		m.cursorPos = m.cursorPos - m.height + 1
	} else {
		m.cursorPos = m.viewportPos
	}
	m.checkCursorVsViewport()
	return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
}

func (m *Model) pgdown() tea.Msg {
	if m.cursorPos == m.viewportPos+m.height-1 {
		m.cursorPos = m.cursorPos + m.height - 1
	} else {
		m.cursorPos = m.viewportPos + m.height - 1
	}
	m.checkCursorVsViewport()
	return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
}

func (m *Model) home() tea.Msg {
	m.cursorPos = 0
	m.checkCursorVsViewport()
	return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
}

func (m *Model) end() tea.Msg {
	m.cursorPos = len(m.completeView) - 1
	m.checkCursorVsViewport()
	return global.SelectGroupMsg{Idx: m.lineToGroup[m.cursorPos]}
}
