package applys

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/box"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/common/messages"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
	"github.com/tiramiseb/quickonf/internal/program/specific/apply"
)

type model struct {
	box    tea.Model
	groups []tea.Model

	srcGroups []*instructions.Group

	width int
}

func New() *model {
	return &model{box: box.New("Applies", "Nothing to apply...", nil, true, false)}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch theMsg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = theMsg.Width - 2 // 2 chars for box border
		apply.GroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.GroupWaiting.Copy().Width(m.width),
			group.StatusRunning:   style.GroupRunning.Copy().Width(m.width),
			group.StatusFailed:    style.GroupFail.Copy().Width(m.width),
			group.StatusSucceeded: style.GroupSuccess.Copy().Width(m.width),
		}
		apply.HoveredGroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.HoveredGroupWaiting.Copy().Width(m.width),
			group.StatusRunning:   style.HoveredGroupRunning.Copy().Width(m.width),
			group.StatusFailed:    style.HoveredGroupFail.Copy().Width(m.width),
			group.StatusSucceeded: style.HoveredGroupSuccess.Copy().Width(m.width),
		}
		apply.SelectedGroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.SelectedGroupWaiting.Copy().Width(m.width),
			group.StatusRunning:   style.SelectedGroupRunning.Copy().Width(m.width),
			group.StatusFailed:    style.SelectedGroupFail.Copy().Width(m.width),
			group.StatusSucceeded: style.SelectedGroupSuccess.Copy().Width(m.width),
		}
	case group.Msg:
		switch theMsg.Type {
		case group.CheckDone:
			cmd, updated := m.maybeAddGroupToApplys(theMsg.Group)
			cmds = append(cmds, cmd)
			if updated {
				m.box, cmd = m.box.Update(box.UpdateElementsMsg{Elements: m.groups})
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		case group.ApplyChange, group.ApplySuccess, group.ApplyFail:
			var cmd tea.Cmd
			m.groups[theMsg.Gidx], cmd = m.groups[theMsg.Gidx].Update(msg)
			cmds = append(cmds, cmd)
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

func (m *model) maybeAddGroupToApplys(grp *instructions.Group) (tea.Cmd, bool) {
	// Check if the group is already in the list of appliable groups
	for i, g := range m.srcGroups {
		if grp == g {
			if len(g.Applys) == 0 {
				m.srcGroups[i] = m.srcGroups[len(m.srcGroups)-1]
				m.srcGroups = m.srcGroups[:len(m.srcGroups)-1]
				return nil, true
			}
			var cmd tea.Cmd
			m.groups[i], cmd = m.groups[i].Update(apply.ResetOutputsMsg{})
			return cmd, true
		}
	}

	// This group is not already in the list, add it
	if len(grp.Applys) > 0 {
		gidx := len(m.groups)
		m.srcGroups = append(m.srcGroups, grp)
		m.groups = append(m.groups, apply.New(grp, gidx, m.width))
		return m.groups[gidx].Init(), true
	}
	return nil, false
}

func (m *model) View() string {
	return m.box.View()
}
