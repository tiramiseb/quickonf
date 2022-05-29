package groups

import (
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

const scrollDelta = 3 // Same value as default value for viewport's mousewheeldelta

func (m *Model) up() {
	m.selectedGroup = m.selectedGroup.Previous(1, !toggles.Get("filter"))
}

func (m *Model) down() {
	m.selectedGroup = m.selectedGroup.Next(1, !toggles.Get("filter"))
}

func (m *Model) pgup() {
	m.selectedGroup = m.selectedGroup.Previous(m.height/2, !toggles.Get("filter"))
}

func (m *Model) pgdown() {
	m.selectedGroup = m.selectedGroup.Next(m.height/2, !toggles.Get("filter"))
}

func (m *Model) home() {
	m.selectedGroup = m.selectedGroup.Previous(m.groups.Count(), !toggles.Get("filter"))
}

func (m *Model) end() {
	m.selectedGroup = m.selectedGroup.Next(m.groups.Count(), !toggles.Get("filter"))
}

func (m *Model) selectLine(lineIdx int) {
	// firstGroupInView := groups.GetSelectedIndex() - m.selectedGroupToViewportOffset
	// clickedGroup := firstGroupInView + lineIdx
	// if clickedGroup >= groups.CountDisplayed() {
	// 	return
	// }
	// groups.IncrementSelected(lineIdx - m.selectedGroupToViewportOffset)
	// m.selectedGroupToViewportOffset = lineIdx
}

func (m *Model) scrollUp() {
	m.selectedGroup = m.selectedGroup.Previous(scrollDelta, !toggles.Get("filter"))
}

func (m *Model) scrollDown() {
	m.selectedGroup = m.selectedGroup.Next(scrollDelta, !toggles.Get("filter"))
}
