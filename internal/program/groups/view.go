package groups

import (
	"strings"

	"github.com/tiramiseb/quickonf/internal/program/global"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

func (m *Model) View() string {
	if m.groups.Count() == 0 {
		return global.MakeWidth("No check", m.width)
	}
	m.firstDisplayedGroup = m.selectedGroup.Previous(m.height/2, !toggles.Get("filter"))
	grp := m.firstDisplayedGroup
	var i int
	var view string
	for i = 0; i < m.height; i++ {
		if grp == m.selectedGroup {
			view += global.SelectedStyles[grp.Status()].Render(global.MakeWidth(grp.Name, m.width)) + "\n"
		} else {
			view += global.Styles[grp.Status()].Render(global.MakeWidth(grp.Name, m.width)) + "\n"
		}
		newGrp := grp.Next(1, !toggles.Get("filter"))
		if newGrp == grp {
			break
		}
		grp = newGrp
	}
	grp = m.firstDisplayedGroup
	for ; i < m.height-1; i++ {
		// There is still space at the bottom, try to add more on top
		newGrp := grp.Previous(1, !toggles.Get("filter"))
		if newGrp == grp {
			break
		}
		grp = newGrp
		view = global.Styles[grp.Status()].Render(global.MakeWidth(grp.Name, m.width)) + "\n" + view
	}
	return strings.TrimRight(view, "\n")
}
