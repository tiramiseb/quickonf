package group

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/state"
)

type status int

type ChangeMessage struct {
	Gidx int
}
type RunningMessage struct {
	Gidx int
}
type FailedMessage struct {
	Gidx int
}
type SucceededMessage struct {
	Gidx int
}

const (
	statusWaiting status = iota
	statusRunning
	statusFailed
	statusSucceeded
)

type Model struct {
	g     *state.Group
	width int

	messages chan interface{}

	idx int

	status status

	outputs []*instructionOutput

	View []string
}

func New(width int, group *state.Group, idx int) *Model {
	return &Model{
		g: group, width: width,
		messages: make(chan interface{}),
		idx:      idx,
		status:   statusWaiting,
	}
}

func (m *Model) listen() tea.Msg {
	return <-m.messages
}

func (m *Model) Run(options state.Options) tea.Cmd {
	go func() {
		m.messages <- RunningMessage{m.idx}
		vars := state.NewVariablesSet()
		for _, ins := range m.g.Instructions {
			if !ins.Run(newInstructionOutput(m, ins.Name()), vars, options) {
				m.messages <- FailedMessage{m.idx}
				return
			}
		}
		m.messages <- SucceededMessage{m.idx}
	}()
	return m.listen
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case ChangeMessage:
		if msg.Gidx != m.idx {
			return nil
		}
		cmd = m.listen
	case RunningMessage:
		if msg.Gidx != m.idx {
			return nil
		}
		m.status = statusRunning
		cmd = m.listen
	case SucceededMessage:
		if msg.Gidx != m.idx {
			return nil
		}
		m.status = statusSucceeded
	case FailedMessage:
		if msg.Gidx != m.idx {
			return nil
		}
		m.status = statusFailed
	}
	m.update()
	return cmd
}

func (m *Model) update() {
	var lines []string
	lines = append(
		lines,
		styleMap[m.status].Copy().Width(m.width).Render(m.g.Name),
	)
	for _, ins := range m.outputs {
		lines = append(
			lines,
			fmt.Sprintf(
				"%s%s",
				instructionStyleMap[ins.Status].Render(ins.name),
				instructionMessageStyle.Render(ins.message),
			),
		)
	}
	m.View = lines
}
