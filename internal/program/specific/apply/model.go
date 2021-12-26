package apply

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

type ResetOutputsMsg struct{}

var GroupStyles = map[group.Status]lipgloss.Style{
	group.StatusWaiting:   style.GroupWaiting,
	group.StatusRunning:   style.GroupRunning,
	group.StatusFailed:    style.GroupFail,
	group.StatusSucceeded: style.GroupSuccess,
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
	status        group.Status
	collapsedView string
	fullView      string
	collapsed     bool

	outputs  []*commandOutput
	messages chan group.Msg
}

func New(grp *instructions.Group, idx, width int) *model {
	m := &model{
		group: grp,
		idx:   idx,

		width: width,

		messages: make(chan group.Msg),
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
	if m.status != group.StatusWaiting {
		return nil
	}
	m.status = group.StatusRunning
	var result group.MsgType
	if m.group.Apply(m.commandOutputs()) {
		result = group.ApplySuccess
	} else {
		result = group.ApplyFail
	}
	return group.Msg{
		Gidx:  m.idx,
		Group: m.group,
		Type:  result,
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
			if m.status == group.StatusWaiting {
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
	case group.Msg:
		if msg.Gidx != m.idx {
			return m, nil
		}
		switch msg.Type {
		case group.ApplyChange:
			cmd = m.listen
		case group.ApplySuccess:
			m.status = group.StatusSucceeded
		case group.ApplyFail:
			m.status = group.StatusFailed
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
