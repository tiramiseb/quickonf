package global

import (
	"github.com/tiramiseb/quickonf/internal/instructions"
)

var (
	AllGroups       []*instructions.Group
	DisplayedGroups []*instructions.Group
	SelectedGroup   int // Index in the "DisplayedGroups" list
)

func GetSelectedGroup() *instructions.Group {
	if len(DisplayedGroups) == 0 {
		return nil
	}
	return DisplayedGroups[SelectedGroup]
}

func init() {
	TogglesListeners["filter"] = append(TogglesListeners["filter"], allGroupsToDisplayedGroups)
}

func GroupsMayHaveChanged() {
	allGroupsToDisplayedGroups(Toggles["filter"])
}

func allGroupsToDisplayedGroups(filtered bool) {
	selectedGroupAddr := DisplayedGroups[SelectedGroup]
	DisplayedGroups = nil
	if filtered {
		i := 0
		for j, g := range AllGroups {
			// TODO If selected group is not displayed, change selected group
			if g.HasApply() {
				DisplayedGroups = append(DisplayedGroups, g)
				if g == selectedGroupAddr {
					SelectedGroup = i
				}
				i++
			} else if g == selectedGroupAddr {
				// If selected group is not displayed, try to select the next one
				if j >= len(AllGroups)-1 {
					SelectedGroup = len(DisplayedGroups) - 1
				} else {
					selectedGroupAddr = AllGroups[j+1]
				}
			}
		}
	} else {
		DisplayedGroups = append(DisplayedGroups, AllGroups...)
		for i, g := range DisplayedGroups {
			if g == selectedGroupAddr {
				SelectedGroup = i
			}
		}
	}
}
