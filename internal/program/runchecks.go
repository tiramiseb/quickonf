package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/instructions"
)

type checkDone struct {
	groupIndex int
}

func orderChecksByPriority(groups []*instructions.Group) [][]int {
	var (
		currentPriority int
		byPriority      [][]int
		thisPriority    []int
	)
	for i, g := range groups {
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
		cmds[i] = m.run(id)
	}
	m.nextPriorityGroup++
	m.currentlyRunningChecks = nbChecks
	return tea.Batch(cmds...)
}

func (m *model) run(i int) tea.Cmd {
	return func() tea.Msg {
		m.groups[i].Run()
		return checkDone{i}
	}
}
