package groups

import (
	"strings"

	"github.com/tiramiseb/quickonf/program/styles"
)

func (m *Model) View() string {
	if m.groups.Count() == 0 {
		return styles.MakeWidth("No check", m.width)
	}
	m.firstDisplayedGroup = m.selectedGroup.Previous(m.height/2, m.showSuccessful)
	grp := m.firstDisplayedGroup
	var i int
	var view string
	for i = 0; i < m.height; i++ {
		if grp == m.selectedGroup {
			view += styles.SelectedStyles[grp.Status()].Render(styles.MakeWidth(grp.Name, m.width)) + "\n"
		} else {
			view += styles.Styles[grp.Status()].Render(styles.MakeWidth(grp.Name, m.width)) + "\n"
		}
		newGrp := grp.Next(1, m.showSuccessful)
		if newGrp == grp {
			break
		}
		grp = newGrp
	}
	grp = m.firstDisplayedGroup
	for ; i < m.height-1; i++ {
		// There is still space at the bottom, try to add more on top
		newGrp := grp.Previous(1, m.showSuccessful)
		if newGrp == grp {
			break
		}
		grp = newGrp
		view = styles.Styles[grp.Status()].Render(styles.MakeWidth(grp.Name, m.width)) + "\n" + view
	}
	return strings.TrimRight(view, "\n")
}
