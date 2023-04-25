package groups

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/instructions"
	"github.com/tiramiseb/quickonf/program/details"
	"github.com/tiramiseb/quickonf/program/messages"
)

type Model struct {
	groups  *instructions.Groups
	details *details.Model

	firstDisplayedGroup *instructions.Group
	selectedGroup       *instructions.Group

	showSuccessful bool

	width  int
	height int
}

func New(d *details.Model) *Model {
	return &Model{
		details:        d,
		showSuccessful: true,
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.up()
			m.provideGroupToDetails()
		case "down":
			m.down()
			m.provideGroupToDetails()
		case "pgup":
			m.pgup()
			m.provideGroupToDetails()
		case "pgdown":
			m.pgdown()
			m.provideGroupToDetails()
		case "home":
			m.home()
			m.provideGroupToDetails()
		case "end":
			m.end()
			m.provideGroupToDetails()
		}
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.scrollUp()
			m.provideGroupToDetails()
		case tea.MouseWheelDown:
			m.scrollDown()
			m.provideGroupToDetails()
		case tea.MouseRelease:
			if msg.Y >= 0 {
				m.selectLine(msg.Y)
				m.provideGroupToDetails()
			}
		}
	case messages.NewSignal:
		m.selectedGroup = m.selectedGroup.Next(0, m.showSuccessful)
		m.provideGroupToDetails()
	}
	return m, nil
}

func (m *Model) provideGroupToDetails() {
	m.details.ShowGroup(m.selectedGroup)
}

func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.width = size.Width
	m.height = size.Height
	return m
}

func (m *Model) ToggleShowSuccessful() tea.Msg {
	m.showSuccessful = !m.showSuccessful
	return messages.ToggleStatus{Name: "filter", Status: !m.showSuccessful}
}

func (m *Model) ApplySelected() tea.Msg {
	m.selectedGroup.Apply()
	return nil
}

func (m *Model) RecheckSelected(signalTarget chan bool) tea.Cmd {
	return func() tea.Msg {
		m.selectedGroup.Check(signalTarget, true)
		return nil
	}
}

func (m *Model) ReplaceGroups(g *instructions.Groups) {
	if m.selectedGroup != nil {
		// Try to keep the same group selected
		grp := g.FirstGroup()
		for {
			if grp.Name == m.selectedGroup.Name {
				// Selected group found, select it
				m.selectedGroup = grp
				break
			}
			newGrp := grp.Next(1, m.showSuccessful)
			if newGrp == grp {
				// Selected group disappeared, select the first one
				m.selectedGroup = g.FirstGroup()
				break
			}
			grp = newGrp
		}
	} else {
		// No group was selected, select the first one
		m.selectedGroup = g.FirstGroup()
	}
	m.groups = g
}
