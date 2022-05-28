package groups

import (
	"github.com/tiramiseb/quickonf/internal/program/global"
	"github.com/tiramiseb/quickonf/internal/program/global/groups"
)

func (m *Model) View() string {
	if groups.CountAll() == 0 {
		return global.MakeWidth("No check", m.width)
	}
	displayed := groups.GetDisplayed()
	countDisplayed := len(displayed)
	if countDisplayed == 0 {
		return global.MakeWidth("No change needed", m.width)
	}
	selected := groups.GetSelectedIndex()
	var view string
	firstGroupInView := selected - m.selectedGroupToViewportOffset
	if firstGroupInView < 0 {
		firstGroupInView = 0
		m.selectedGroupToViewportOffset = selected
	}
	afterLastGroupInView := firstGroupInView + m.height
	freeLinesAtBottom := afterLastGroupInView - countDisplayed
	if freeLinesAtBottom > 0 {
		afterLastGroupInView = countDisplayed
		firstGroupInView -= freeLinesAtBottom
		if firstGroupInView < 0 {
			firstGroupInView = 0
			m.selectedGroupToViewportOffset = selected
		} else {
			m.selectedGroupToViewportOffset = selected - firstGroupInView
		}
	}
	for i := firstGroupInView; i < afterLastGroupInView-1; i++ {
		g := displayed[i]
		if i == selected {
			view += global.SelectedStyles[g.Status()].Render(global.MakeWidth(g.Name, m.width)) + "\n"
		} else {
			view += global.Styles[g.Status()].Render(global.MakeWidth(g.Name, m.width)) + "\n"
		}
	}
	lastGroup := afterLastGroupInView - 1
	if lastGroup >= 0 {
		g := displayed[lastGroup]
		if lastGroup == selected {
			view += global.SelectedStyles[g.Status()].Render(global.MakeWidth(g.Name, m.width))
		} else {
			view += global.Styles[g.Status()].Render(global.MakeWidth(g.Name, m.width))
		}
	}

	return view
}
