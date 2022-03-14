package checks

import (
	"github.com/tiramiseb/quickonf/internal/program/global"
)

const scrollDelta = 3 // Same value as default value for viewport's mousewheeldelta

func (m *Model) centerOnSelectedGroup() {
	m.selectedGroupToViewportOffset = m.height / 2
}

func (m *Model) up() {
	if global.SelectedGroup == 0 {
		return
	}
	global.SelectedGroup--
	m.centerOnSelectedGroup()
}

func (m *Model) down() {
	if global.SelectedGroup == len(global.DisplayedGroups)-1 {
		// Already last group
		return
	}
	global.SelectedGroup++
	m.centerOnSelectedGroup()
}

func (m *Model) pgup() {
	global.SelectedGroup -= m.height / 2
	if global.SelectedGroup < 0 {
		global.SelectedGroup = 0
	}
	m.centerOnSelectedGroup()
}

func (m *Model) pgdown() {
	global.SelectedGroup += m.height / 2
	if global.SelectedGroup >= len(global.DisplayedGroups) {
		global.SelectedGroup = len(global.DisplayedGroups) - 1
	}
	m.centerOnSelectedGroup()
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

func (m *Model) selectLine(lineIdx int) {
	global.SelectedGroup = global.SelectedGroup - m.selectedGroupToViewportOffset + lineIdx
	m.selectedGroupToViewportOffset = lineIdx
}

func (m *Model) scrollUp() {
	if global.SelectedGroup-m.selectedGroupToViewportOffset <= 0 {
		// Already at top of view, do not scroll more
		return
	}
	m.selectedGroupToViewportOffset += scrollDelta
}

func (m *Model) scrollDown() {
	if global.SelectedGroup-m.selectedGroupToViewportOffset+m.height >= len(global.DisplayedGroups) {
		// Already at bottom of view, do not scroll more
		return
	}
	m.selectedGroupToViewportOffset -= scrollDelta
}
