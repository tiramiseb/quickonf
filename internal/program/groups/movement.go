package groups

import (
	"github.com/tiramiseb/quickonf/internal/program/global/groups"
)

const scrollDelta = 3 // Same value as default value for viewport's mousewheeldelta

func (m *Model) centerOnSelectedGroup() {
	m.selectedGroupToViewportOffset = m.height / 2
}

func (m *Model) up() {
	groups.DecrementSelected(1)
	m.centerOnSelectedGroup()
}

func (m *Model) down() {
	groups.IncrementSelected(1)
	m.centerOnSelectedGroup()
}

func (m *Model) pgup() {
	groups.DecrementSelected(m.height / 2)
	m.centerOnSelectedGroup()
}

func (m *Model) pgdown() {
	groups.IncrementSelected(m.height / 2)
	m.centerOnSelectedGroup()
}

func (m *Model) home() {
	groups.SelectFirst()
	m.selectedGroupToViewportOffset = 0
}

func (m *Model) end() {
	groups.SelectLast()
	m.selectedGroupToViewportOffset = groups.CountDisplayed() - m.height
	if m.selectedGroupToViewportOffset < 0 {
		// Less groups than available space, display all
		m.selectedGroupToViewportOffset = groups.CountDisplayed() - 1
	}
}

func (m *Model) selectLine(lineIdx int) {
	firstGroupInView := groups.GetSelectedIndex() - m.selectedGroupToViewportOffset
	clickedGroup := firstGroupInView + lineIdx
	if clickedGroup >= groups.CountDisplayed() {
		return
	}
	groups.IncrementSelected(lineIdx - m.selectedGroupToViewportOffset)
	m.selectedGroupToViewportOffset = lineIdx
}

func (m *Model) scrollUp() {
	if groups.GetSelectedIndex()-m.selectedGroupToViewportOffset <= 0 {
		// Already at top of view, do not scroll more
		return
	}
	m.selectedGroupToViewportOffset += scrollDelta
}

func (m *Model) scrollDown() {
	if groups.GetSelectedIndex()-m.selectedGroupToViewportOffset+m.height >= groups.CountDisplayed() {
		// Already at bottom of view, do not scroll more
		return
	}
	m.selectedGroupToViewportOffset -= scrollDelta
}
