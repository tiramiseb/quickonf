package checks

import (
	"github.com/tiramiseb/quickonf/internal/program/global"
)

func (m *Model) View() string {
	if len(global.AllGroups) == 0 {
		return global.MakeWidth("No check", m.width)
	}
	if len(global.DisplayedGroups) == 0 {
		return global.MakeWidth("No change needed", m.width)
	}
	var view string
	firstGroupInView := global.SelectedGroup - m.selectedGroupToViewportOffset
	if firstGroupInView < 0 {
		firstGroupInView = 0
		m.selectedGroupToViewportOffset = global.SelectedGroup
	}
	afterLastGroupInView := firstGroupInView + m.height
	freeLinesAtBottom := afterLastGroupInView - len(global.DisplayedGroups)
	if freeLinesAtBottom > 0 {
		afterLastGroupInView = len(global.DisplayedGroups)
		firstGroupInView -= freeLinesAtBottom
		if firstGroupInView < 0 {
			firstGroupInView = 0
			m.selectedGroupToViewportOffset = global.SelectedGroup
		}
	}
	for i := firstGroupInView; i < afterLastGroupInView-1; i++ {
		g := global.DisplayedGroups[i]
		if i == global.SelectedGroup {
			view += global.SelectedStyles[g.Status()].Render(global.MakeWidth(g.Name, m.width)) + "\n"
		} else {
			view += global.Styles[g.Status()].Render(global.MakeWidth(g.Name, m.width)) + "\n"
		}
	}
	lastGroup := afterLastGroupInView - 1
	if lastGroup >= 0 {
		g := global.DisplayedGroups[lastGroup]
		if lastGroup == global.SelectedGroup {
			view += global.SelectedStyles[g.Status()].Render(global.MakeWidth(g.Name, m.width))
		} else {
			view += global.Styles[g.Status()].Render(global.MakeWidth(g.Name, m.width))
		}
	}

	return view
}
