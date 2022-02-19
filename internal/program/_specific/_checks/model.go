package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/box"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/common/messages"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
	"github.com/tiramiseb/quickonf/internal/program/specific/check"
)

type model struct {
	box    tea.Model
	groups []tea.Model

	byPriority             [][]int
	nextPriorityGroup      int
	currentlyRunningChecks int
}

func New(groups []*instructions.Group) *model {
	gs := make([]tea.Model, len(groups))
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
		gs[i] = check.New(i, g)
	}
	if thisPriority != nil {
		byPriority = append(byPriority, thisPriority)
	}
	return &model{
		box:        box.New("Checks", "Nothing to check...", gs, false, true),
		groups:     gs,
		byPriority: byPriority,
	}
}

func (m *model) Init() tea.Cmd {
	return m.next()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch theMsg := msg.(type) {
	case tea.WindowSizeMsg:
		groupWidth := theMsg.Width - 2 // 2 chars for box border
		check.GroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.GroupWaiting.Copy().Width(groupWidth),
			group.StatusInfo:      style.GroupInfo.Copy().Width(groupWidth),
			group.StatusRunning:   style.GroupRunning.Copy().Width(groupWidth),
			group.StatusFailed:    style.GroupFail.Copy().Width(groupWidth),
			group.StatusSucceeded: style.GroupSuccess.Copy().Width(groupWidth),
		}
		check.HoveredGroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.HoveredGroupWaiting.Copy().Width(groupWidth),
			group.StatusInfo:      style.HoveredGroupInfo.Copy().Width(groupWidth),
			group.StatusRunning:   style.HoveredGroupRunning.Copy().Width(groupWidth),
			group.StatusFailed:    style.HoveredGroupFail.Copy().Width(groupWidth),
			group.StatusSucceeded: style.HoveredGroupSuccess.Copy().Width(groupWidth),
		}
		check.SelectedGroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.SelectedGroupWaiting.Copy().Width(groupWidth),
			group.StatusInfo:      style.SelectedGroupInfo.Copy().Width(groupWidth),
			group.StatusRunning:   style.SelectedGroupRunning.Copy().Width(groupWidth),
			group.StatusFailed:    style.SelectedGroupFail.Copy().Width(groupWidth),
			group.StatusSucceeded: style.SelectedGroupSuccess.Copy().Width(groupWidth),
		}
	case group.Msg:
		switch theMsg.Type {
		case group.CheckTrigger:
			var cmd tea.Cmd
			m.groups[theMsg.Gidx], cmd = m.groups[theMsg.Gidx].Update(msg)
			cmds = append(cmds, cmd)
		case group.CheckDone:
			m.currentlyRunningChecks--
			var cmd tea.Cmd
			m.groups[theMsg.Gidx], cmd = m.groups[theMsg.Gidx].Update(msg)
			cmds = append(cmds, cmd)
			if m.currentlyRunningChecks == 0 {
				cmds = append(cmds, m.next())
			}
		}
		msg = box.ForceRedrawMsg{}
	case messages.FilterMsg:
		for i, g := range m.groups {
			var cmd tea.Cmd
			m.groups[i], cmd = g.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		msg = box.ForceRedrawMsg{}
	}
	var cmd tea.Cmd
	m.box, cmd = m.box.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return m.box.View()
}
