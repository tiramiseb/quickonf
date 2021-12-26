package checks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/box"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
	"github.com/tiramiseb/quickonf/internal/program/specific/check"
)

type model struct {
	box    tea.Model
	groups []tea.Model

	nextInQueue int
}

func New(groups []*instructions.Group) *model {
	gs := make([]tea.Model, len(groups))
	for i, g := range groups {
		gs[i] = check.New(i, g)
	}
	return &model{
		box:    box.New("Checks", "Nothing to check...", gs, false, true),
		groups: gs,
	}
}

func (m *model) Init() tea.Cmd {
	// Need to discriminate on priority before running in parallel
	// nb := runtime.NumCPU()
	// cmds := make([]tea.Cmd, nb)
	// for i := 0; i < nb; i++ {
	// 	cmds[i] = m.next()
	// }
	// return tea.Batch(cmds...)
	return m.next()
}

func (m *model) next() tea.Cmd {
	if m.nextInQueue >= len(m.groups) {
		return nil
	}
	cmd := m.groups[m.nextInQueue].Init()
	m.nextInQueue++
	return cmd
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd1 tea.Cmd
		cmd2 tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		groupWidth := msg.Width - 2 // 2 chars for box border
		check.GroupStyles = map[group.Status]lipgloss.Style{
			group.StatusWaiting:   style.GroupWaiting.Copy().Width(groupWidth),
			group.StatusRunning:   style.GroupRunning.Copy().Width(groupWidth),
			group.StatusFailed:    style.GroupFail.Copy().Width(groupWidth),
			group.StatusSucceeded: style.GroupSuccess.Copy().Width(groupWidth),
		}
	case group.Msg:
		switch msg.Type {
		case group.CheckTrigger:
			m.groups[msg.Gidx], cmd1 = m.groups[msg.Gidx].Update(msg)
		case group.CheckDone:
			cmds := make([]tea.Cmd, 2)
			m.groups[msg.Gidx], cmds[0] = m.groups[msg.Gidx].Update(msg)
			cmds[1] = m.next()
			cmd1 = tea.Batch(cmds...)
		}
	}
	m.box, cmd2 = m.box.Update(msg)
	return m, tea.Batch(cmd1, cmd2)
}

func (m *model) View() string {
	return m.box.View()
}
