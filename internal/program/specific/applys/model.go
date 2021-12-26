package applys

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/box"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
	"github.com/tiramiseb/quickonf/internal/program/specific/apply"
)

type groupLine struct {
	gidx      int
	groupline int
}

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
	var (
		cmd1 tea.Cmd
		cmd2 tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2 // 2 chars for box border
		apply.GroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.GroupWaiting.Copy().Width(m.width),
			group.StatusRunning:   style.GroupRunning.Copy().Width(m.width),
			group.StatusFailed:    style.GroupFail.Copy().Width(m.width),
			group.StatusSucceeded: style.GroupSuccess.Copy().Width(m.width),
		}
	case group.Msg:
		switch msg.Type {
		case group.CheckDone:
			var updated bool
			cmd1, updated = m.maybeAddGroupToApplys(msg.Group)
			if updated {
				m.box, cmd2 = m.box.Update(box.UpdateGroupsMsg{Groups: m.groups})
				return m, tea.Batch(cmd1, cmd2)
			}
		case group.ApplyChange, group.ApplySuccess, group.ApplyFail:
			m.groups[msg.Gidx], cmd1 = m.groups[msg.Gidx].Update(msg)
		}
	}
	m.box, cmd2 = m.box.Update(msg)
	return m, tea.Batch(cmd1, cmd2)
}

func (m *model) maybeAddGroupToApplys(grp *instructions.Group) (tea.Cmd, bool) {
	alreadyIncluded := false

	// Check if the group is already in the list of appliable groups
	for i, g := range m.srcGroups {
		if grp == g {
			alreadyIncluded = true
			if len(g.Applys) == 0 {
				m.srcGroups[i] = m.srcGroups[len(m.srcGroups)-1]
				m.srcGroups = m.srcGroups[:len(m.srcGroups)-1]
			} else {
				g.Reports = nil
			}
		}
	}

	// This group is not already in the list, add it
	if !alreadyIncluded && len(grp.Applys) > 0 {
		gidx := len(m.groups)
		m.srcGroups = append(m.srcGroups, grp)
		m.groups = append(m.groups, apply.New(grp, gidx, m.width-2))
		return m.groups[gidx].Init(), true
	}
	return nil, false
}

func (m *model) View() string {
	return m.box.View()
}
