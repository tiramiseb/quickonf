package apply

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

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

type SuccessMsg struct {
	Gidx int
}

type FailMsg struct {
	Gidx int
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

	outputs  []*commandOutput
	messages chan ChangeMsg
}

func New(group *instructions.Group, idx, width int) *model {
	m := &model{
		group: group,
		idx:   idx,

		width: width,

		messages: make(chan ChangeMsg),
	}
	m.updateGroupname()
	m.updateView()
	return m
}

func (m *model) listen() tea.Msg {
	return <-m.messages
}

func (m *model) Init() tea.Cmd {
	return m.listen
}

func (m *model) run() tea.Msg {
	// TODO Allow re-running...
	if m.status != StatusWaiting {
		return nil
	}
	if m.group.Apply(m.commandOutputs()) {
		return SuccessMsg{m.idx}
	} else {
		return FailMsg{m.idx}
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
		case " ", "t", "T":
			m.collapsed = !m.collapsed
		case "enter", "x", "X":
			if m.status == StatusWaiting {
				return m, m.run
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
	case ChangeMsg:
		if msg.Gidx != m.idx {
			return m, nil
		}
		cmd = m.listen
	case SuccessMsg:
		if msg.Gidx != m.idx {
			return m, nil
		}
		m.status = StatusSucceeded
	case FailMsg:
		if msg.Gidx != m.idx {
			return m, nil
		}
		m.status = StatusFailed
	}
	m.updateView()
	return m, cmd
}

func (m *model) updateView() {
	m.collapsedView = GroupStyles[m.status].Render("⏵ " + m.groupName)
	lines := []string{
		GroupStyles[m.status].Render("⏷ " + m.groupName),
	}
	if len(m.outputs) == 0 {
		for _, apply := range m.group.Applys {
			lines = append(
				lines,
				m.instructionLine(
					apply.Name,
					commands.StatusInfo,
					apply.Intro,
				),
			)
		}
	} else {
		for _, message := range m.outputs {
			lines = append(
				lines,
				m.instructionLine(
					message.name,
					message.status,
					message.message,
				),
			)
		}
	}
	m.fullView = strings.Join(lines, "\n")
}

func (m *model) View() string {
	if m.collapsed {
		return m.collapsedView
	}
	return m.fullView
}
