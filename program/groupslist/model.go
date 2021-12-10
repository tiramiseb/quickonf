package groupslist

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/program/group"
	"github.com/tiramiseb/quickonf/state"
)

const moveStepDuration = 15 * time.Millisecond

type moveUpMessage struct{}

type moveDownMessage struct{}

type Model struct {
	width           int
	height          int
	verticalMargins int

	nextInQueue int // index of the next group to run

	keepLines    int // number of lines to keep in the "current group", while moving
	moving       bool
	currentGroup int // index of the group to display on the bottom of the screen - when <0, display the last running group

	groups []*group.Model
	state  *state.State
}

func New(width, height, verticalMargins int, st *state.State) *Model {
	groups := make([]*group.Model, len(st.Groups))
	for i, g := range st.Groups {
		groups[i] = group.New(width, g, i)
	}
	return &Model{
		width,
		height,
		verticalMargins,

		0,

		0,
		false,
		-1,

		groups,
		st,
	}
}

func (m *Model) runNext() tea.Cmd {
	if m.nextInQueue >= len(m.groups) {
		return nil
	}
	cmd := m.groups[m.nextInQueue].Run(m.state.Options)
	m.nextInQueue++
	return cmd
}

func (m *Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, m.state.Options.NbConcurrentGroups)
	for i := 0; i < m.state.Options.NbConcurrentGroups; i++ {
		cmds[i] = m.runNext()
	}
	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			cmd = m.startMoveUp
		case "down":
			cmd = m.startMoveDown
		case "right":
			m.currentGroup = -1
		}
	case moveUpMessage:
		cmd = m.waitMoveUp
	case moveDownMessage:
		cmd = m.waitMoveDown
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - m.verticalMargins
		for i := range m.groups {
			m.groups[i], _ = m.groups[i].Update(msg)
		}
	case group.ChangeMessage:
		m.groups[msg.Gidx], cmd = m.groups[msg.Gidx].Update(msg)
	case group.RunningMessage:
		m.groups[msg.Gidx], cmd = m.groups[msg.Gidx].Update(msg)
	case group.SucceededMessage:
		m.groups[msg.Gidx], _ = m.groups[msg.Gidx].Update(msg)
		cmd = m.runNext()
	case group.FailedMessage:
		m.groups[msg.Gidx], _ = m.groups[msg.Gidx].Update(msg)
		cmd = m.runNext()
	}
	return cmd
}

func (m *Model) waitMoveUp() tea.Msg {
	time.Sleep(moveStepDuration)
	return m.moveUp()
}

func (m *Model) waitMoveDown() tea.Msg {
	time.Sleep(moveStepDuration)
	return m.moveDown()
}

func (m *Model) startMoveUp() tea.Msg {
	if m.moving {
		return nil
	}
	if m.currentGroup == -1 {
		m.currentGroup = m.nextInQueue - 1
	}
	if m.currentGroup == 0 {
		return nil
	}
	m.moving = true
	m.keepLines = len(m.groups[m.currentGroup].View)
	return m.moveUp()
}

func (m *Model) startMoveDown() tea.Msg {
	if m.moving {
		return nil
	}
	if m.currentGroup == -1 {
		m.currentGroup = m.nextInQueue - 1
	}
	if m.currentGroup == len(m.groups)-1 {
		return nil
	}
	m.currentGroup++
	m.moving = true
	m.keepLines = 0
	return m.moveDown()
}

func (m *Model) moveUp() tea.Msg {
	m.keepLines--
	if m.keepLines == 0 {
		m.moving = false
		m.currentGroup--
		return nil
	}
	return moveUpMessage{}
}

func (m *Model) moveDown() tea.Msg {
	m.keepLines++
	if m.keepLines == len(m.groups[m.currentGroup].View) {
		m.moving = false
		m.keepLines = 0
		return nil
	}
	return moveDownMessage{}
}

func (m *Model) View() string {
	var views []string
	var displayedHeight int
	currentGroup := m.currentGroup
	if m.currentGroup == -1 {
		currentGroup = m.nextInQueue - 1
	}
	for i := currentGroup; i >= 0; i-- {
		groupView := m.groups[i].View
		if i == currentGroup && m.keepLines != 0 {
			groupView = groupView[:m.keepLines]
		}
		displayedHeight += len(groupView)
		if displayedHeight > m.height {
			delta := displayedHeight - m.height
			groupView = groupView[delta:]
			views = append(groupView, views...)
			break
		}
		views = append(groupView, views...)
	}
	if displayedHeight < m.height {
		views = append([]string{strings.Repeat("\n", m.height-displayedHeight-1)}, views...)
	}
	return strings.Join(views, "\n")
}
