package check

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

type TriggerMsg struct {
	Gidx int
}

type DoneMsg struct {
	Gidx  int
	Group *instructions.Group
}

type Status int

const (
	StatusWaiting Status = iota
	StatusRunning
	StatusFailed
	StatusSucceeded
)

var GroupStyles = map[Status]lipgloss.Style{
	StatusWaiting:   style.GroupWaiting,
	StatusRunning:   style.GroupRunning,
	StatusFailed:    style.GroupFail,
	StatusSucceeded: style.GroupSuccess,
}

var InstructionStyles = map[commands.Status]lipgloss.Style{
	commands.StatusInfo:    style.InstructionInfo,
	commands.StatusError:   style.InstructionError,
	commands.StatusSuccess: style.InstructionSuccess,
}

type model struct {
	group *instructions.Group
	idx   int

	width         int
	groupName     string
	status        Status
	collapsedView string
	fullView      string
	collapsed     bool
}

func New(i int, g *instructions.Group) *model {
	return &model{
		group:     g,
		idx:       i,
		width:     2,
		status:    StatusWaiting,
		collapsed: true,
	}
}

func (m *model) Init() tea.Cmd {
	return m.trigger
}

func (m *model) trigger() tea.Msg {
	if m.status == StatusWaiting {
		m.status = StatusRunning
		return TriggerMsg{m.idx}
	}
	return nil
}

func (m *model) run() tea.Msg {
	if m.group.Run() {
		m.status = StatusSucceeded
	} else {
		m.status = StatusFailed
		m.collapsed = false
	}
	return DoneMsg{m.idx, m.group}
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.updateGroupname()
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "t", "T":
			m.collapsed = !m.collapsed
		case "enter", "x", "X":
			if m.status != StatusRunning {
				m.group.Reset()
				m.status = StatusWaiting
				return m, m.trigger
			}
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseUnknown:
			return m, nil
		case tea.MouseRelease:
			if msg.Y == 0 {
				m.collapsed = !m.collapsed
			}
		}
	case TriggerMsg:
		if msg.Gidx != m.idx {
			return m, nil
		}
		cmd = m.run
	case DoneMsg:
		if msg.Gidx != m.idx {
			return m, nil
		}
	}
	m.updateView()
	return m, cmd
}

func (m *model) updateView() {
	m.collapsedView = GroupStyles[m.status].Render("⏵ " + m.groupName)
	lines := []string{
		GroupStyles[m.status].Render("⏷ " + m.groupName),
	}
	if len(m.group.Reports) > 0 {
		for _, report := range m.group.Reports {
			lines = append(
				lines,
				m.instructionLine(report),
			)
		}
	} else {
		lines = append(
			lines,
			m.instructionLine(instructions.CheckReport{
				Name:   "Empty",
				Status: commands.StatusInfo,
			}),
		)
	}
	m.fullView = strings.Join(lines, "\n")
}

func (m *model) View() string {
	if m.collapsed {
		return m.collapsedView
	}
	return m.fullView
}
