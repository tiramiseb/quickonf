package checks

import (
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) up() {
	if global.SelectedGroup == 0 {
		// Already first group
		return
	}
	global.SelectedGroup--
	m.selectedGroupToViewportOffset--
	if m.selectedGroupToViewportOffset < 0 {
		m.selectedGroupToViewportOffset = 0
	}
}

func (m *Model) down() {
	if global.SelectedGroup == len(global.DisplayedGroups)-1 {
		// Already last group
		return
	}
	global.SelectedGroup++
	m.selectedGroupToViewportOffset++
	if m.selectedGroupToViewportOffset >= m.height {
		m.selectedGroupToViewportOffset = m.height - 1
	}
}

func (m *Model) pgup() {
	switch {
	case global.SelectedGroup == 0:
		// Already first group
	case m.selectedGroupToViewportOffset == 0:
		// Selected group is first line, jump 1 screen
		global.SelectedGroup -= m.height
		if global.SelectedGroup < 0 {
			global.SelectedGroup = 0
		}
	default:
		// Selected group is not first line, select first line
		global.SelectedGroup -= m.selectedGroupToViewportOffset
		m.selectedGroupToViewportOffset = 0
	}
}

func (m *Model) pgdown() {
	switch {
	case global.SelectedGroup == len(global.DisplayedGroups)-1:
		// Already last group
	case m.selectedGroupToViewportOffset == m.height-1:
		// Selected group is last line, jump 1 screen
		global.SelectedGroup += m.height
		if global.SelectedGroup >= len(global.DisplayedGroups) {
			global.SelectedGroup = len(global.DisplayedGroups) - 1
		}
	default:
		// Selected group is not last line, select last line
		global.SelectedGroup += (m.height - m.selectedGroupToViewportOffset - 1)
		m.selectedGroupToViewportOffset = m.height - 1
	}
}

func (m *Model) home() {
	global.SelectedGroup = 0
	m.selectedGroupToViewportOffset = 0
}

func (m *Model) end() {
	global.SelectedGroup = len(global.DisplayedGroups) - 1
	m.selectedGroupToViewportOffset = len(global.DisplayedGroups) - m.height
	if m.selectedGroupToViewportOffset < 0 {
		// Less groups than available space, display all
		m.selectedGroupToViewportOffset = len(global.DisplayedGroups) - 1
	}
}
