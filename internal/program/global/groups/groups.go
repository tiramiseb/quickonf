package groups

import (
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

var (
	selected      int // Index in the "DisplayedGroups" list
	maxNameLength int
)

func Initialize(g []*instructions.Group) {
	all = g
	displayed = g
	selected = 0
	for _, group := range g {
		l := len(group.Name)
		if l > maxNameLength {
			maxNameLength = l
		}
	}
}

func init() {
	toggles.AddListener("filter", allGroupsToDisplayedGroups)
}

func GetMaxNameLength() int {
	return maxNameLength
}

func allGroupsToDisplayedGroups(filtered bool) {
	if len(displayed) == 0 {
		return
	}
	selectedGroup := displayed[selected]
	displayed = nil
	if filtered {
		i := 0
		for j, g := range all {
			if g.Status() != commands.StatusSuccess {
				displayed = append(displayed, g)
				if g == selectedGroup {
					selected = i
				}
				i++
			} else if g == selectedGroup {
				// If selected group is not displayed, try to select the next one
				if j >= len(all)-1 {
					selected = len(displayed) - 1
				} else {
					selectedGroup = all[j+1]
				}
			}
		}
	} else {
		displayed = append(displayed, all...)
		for i, g := range displayed {
			if g == selectedGroup {
				selected = i
			}
		}
	}
}
