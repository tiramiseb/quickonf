package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/global/groups"
)

type checkDone struct{}

func checksIndexByPriority() [][]int {
	var (
		currentPriority int
		byPriority      [][]int
		thisPriority    []int
	)
	for i, g := range groups.GetAll() {
		if g.Priority != currentPriority {
			if thisPriority != nil {
				byPriority = append(byPriority, thisPriority)
			}
			thisPriority = nil
			currentPriority = g.Priority
		}
		thisPriority = append(thisPriority, i)
	}
	if thisPriority != nil {
		byPriority = append(byPriority, thisPriority)
	}
	return byPriority
}

func (m *model) next() tea.Cmd {
	if m.nextPriorityGroup >= len(m.byPriority) {
		return nil
	}
	groupIDs := m.byPriority[m.nextPriorityGroup]
	nbChecks := len(groupIDs)
	cmds := make([]tea.Cmd, nbChecks)
	for i, id := range groupIDs {
		thisIdx := i
		thisId := id
		cmds[thisIdx] = func() tea.Msg {
			groups.InitialCheck(thisId, m.signalTarget)
			return checkDone{}
		}
	}
	m.nextPriorityGroup++
	m.currentlyRunningChecks = nbChecks
	return tea.Batch(cmds...)
}
