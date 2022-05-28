package global

import (
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
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
	toggles.AddListener("filter", allGroupsToDisplayedGroups)
}

func GroupsMayHaveChanged() {
	allGroupsToDisplayedGroups(toggles.Get("filter"))
}

func allGroupsToDisplayedGroups(filtered bool) {
	var selectedGroupAddr *instructions.Group
	if DisplayedGroups == nil {
		selectedGroupAddr = AllGroups[0]
	} else {
		selectedGroupAddr = DisplayedGroups[SelectedGroup]
	}
	DisplayedGroups = nil
	if filtered {
		i := 0
		for j, g := range AllGroups {
			if g.Status() != commands.StatusSuccess {
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
