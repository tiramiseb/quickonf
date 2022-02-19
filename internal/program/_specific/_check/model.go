package check

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/box"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/common/messages"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

var (
	GroupStyles = map[group.Status]lipgloss.Style{
		group.StatusWaiting:   style.GroupWaiting,
		group.StatusInfo:      style.GroupInfo,
		group.StatusRunning:   style.GroupRunning,
		group.StatusFailed:    style.GroupFail,
		group.StatusSucceeded: style.GroupSuccess,
	}
	HoveredGroupStyles = map[group.Status]lipgloss.Style{
		group.StatusWaiting:   style.HoveredGroupWaiting,
		group.StatusInfo:      style.HoveredGroupInfo,
		group.StatusRunning:   style.HoveredGroupRunning,
		group.StatusFailed:    style.HoveredGroupFail,
		group.StatusSucceeded: style.HoveredGroupSuccess,
	}
	SelectedGroupStyles = map[group.Status]lipgloss.Style{
		group.StatusWaiting:   style.SelectedGroupWaiting,
		group.StatusInfo:      style.SelectedGroupInfo,
		group.StatusRunning:   style.SelectedGroupRunning,
		group.StatusFailed:    style.SelectedGroupFail,
		group.StatusSucceeded: style.SelectedGroupSuccess,
	}
)

type model struct {
	group *instructions.Group
	idx   int

	width     int
	groupName string
	status    group.Status
	view      string
	hovered   bool
	selected  bool

	filtered bool
}

func New(i int, g *instructions.Group) *model {
	return &model{
		group:  g,
		idx:    i,
		width:  2,
		status: group.StatusWaiting,
	}
}

func (m *model) Init() tea.Cmd {
	return m.trigger
}

func (m *model) trigger() tea.Msg {
	if m.status == group.StatusWaiting {
		m.status = group.StatusRunning
		return group.Msg{
			Gidx:  m.idx,
			Group: m.group,
			Type:  group.CheckTrigger,
		}
	}
	return nil
}

func (m *model) run() tea.Msg {
	if m.group.Run() {
		if m.group.HasApply() {
			m.status = group.StatusInfo
		} else {
			m.status = group.StatusSucceeded
		}
	} else {
		m.status = group.StatusFailed
	}
	return group.Msg{
		Gidx:  m.idx,
		Group: m.group,
		Type:  group.CheckDone,
	}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.updateGroupname()
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "x", "X":
			if m.status != group.StatusRunning {
				m.group.Reset()
				m.status = group.StatusWaiting
				return m, m.trigger
			}
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseUnknown:
			if !m.hovered {
				return m, nil
			}
			m.hovered = false
		default:
			if !m.hovered {
				m.hovered = true
			}
		}
	case group.Msg:
		if msg.Gidx != m.idx {
			return m, nil
		}
		switch msg.Type {
		case group.CheckTrigger:
			cmd = m.run
		}
	case box.ElementSelectedMsg:
		m.selected = msg.Selected
	case messages.FilterMsg:
		m.filtered = msg.On
	}
	m.updateView()
	return m, cmd
}

func (m *model) updateView() {
	var groupstyle lipgloss.Style
	switch {
	case m.hovered:
		groupstyle = HoveredGroupStyles[m.status]
	case m.selected:
		groupstyle = SelectedGroupStyles[m.status]
	default:
		groupstyle = GroupStyles[m.status]
	}
	m.view = groupstyle.Render(m.groupName)
}

func (m *model) View() string {
	if m.filtered && m.status == group.StatusSucceeded {
		return ""
	}
	return m.view
}
