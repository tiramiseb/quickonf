package groups

import (
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
)

func ApplySelected() {
	if len(displayed) == 0 {
		return
	}
	displayed[selected].Apply()
}

func InitialCheck(idx int, signalTarget chan bool) {
	group := all[idx]
	group.Check(signalTarget)
	if toggles.Get("filter") && group.Status() == commands.StatusSuccess {
		displayedMu.Lock()
		// Filter is enabled and check succeeded, remove this check from displayed groups
		newDisplayed := []*instructions.Group{}
		for _, dispGroup := range displayed {
			if dispGroup != group {
				newDisplayed = append(newDisplayed, dispGroup)
			}
		}
		displayed = newDisplayed
		displayedMu.Unlock()
	}
}
