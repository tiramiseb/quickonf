package applys

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/instructions"
	"github.com/tiramiseb/quickonf/internal/program/apply"
	"github.com/tiramiseb/quickonf/internal/program/check"
	"github.com/tiramiseb/quickonf/internal/program/separator"
)

type ActiveMsg struct {
	Active bool
}

type groupLine struct {
	gidx      int
	groupline int
}

type model struct {
	srcGroups []*instructions.Group
	groups    []tea.Model

	width         int
	boxHeight     int
	subtitleStyle lipgloss.Style
	boxStyle      lipgloss.Style
	active        bool

	allGroupsView  []string    // All groups! And then take a window of this for view
	allLineToGroup []groupLine // map of line number to displayed group, for passing clicks to the correct group

	view            string // The resulting view itself
	viewLineToGroup []groupLine

	selectedGroup          int // index of the selected group in the groups list
	selectedGroupFirstLine int // line on screen of the first line of the selected group (for cursor position)
}

func New() *model {
	return &model{width: 2, boxHeight: 1}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		cmd = m.windowSize(msg)
		m.redrawContent()
	case tea.KeyMsg:
		if len(m.groups) == 0 {
			break
		}
		switch msg.String() {
		case "up":
			m.cursorUp()
		case "down":
			m.cursorDown()
		default:
			m.groups[m.selectedGroup], cmd = m.groups[m.selectedGroup].Update(msg)
			m.redrawContent()
		}
	case tea.MouseMsg:
		// MouseUnknown means the user clicked somewhere else, inform everyone
		cmds := make([]tea.Cmd, len(m.groups))
		if msg.Type == tea.MouseUnknown {
			for i, g := range m.groups {
				m.groups[i], cmds[i] = g.Update(msg)
			}
			cmd = tea.Batch(cmds...)
			m.redrawContent()
			break
		}

		// These situations mean the user did not click on any group
		if msg.X <= 0 || msg.X >= m.width-1 || msg.Y <= 1 || msg.Y >= m.boxHeight+2 {
			break
		}

		// Mouse wheel movement
		if msg.Type == tea.MouseWheelUp {
			m.cursorUp()
			break
		}
		if msg.Type == tea.MouseWheelDown {
			m.cursorDown()
			break
		}

		// Forward mouse message to the group under the cursor
		unknown := tea.MouseMsg{
			Type: tea.MouseUnknown,
		}
		msg.X--
		line := msg.Y - 2
		clickedGroup := -1
		if line < len(m.viewLineToGroup) {
			groupline := m.viewLineToGroup[msg.Y-2]
			msg.Y = groupline.groupline
			m.groups[groupline.gidx], cmds[groupline.gidx] = m.groups[groupline.gidx].Update(msg)
			clickedGroup = groupline.gidx
		}

		for i, g := range m.groups {
			if i == clickedGroup {
				continue
			}
			m.groups[i], cmds[i] = g.Update(unknown)
		}
		cmd = tea.Batch(cmds...)
		m.redrawContent()
	case check.DoneMsg:
		alreadyIncluded := false
		for i, g := range m.srcGroups {
			if msg.Group == g {
				alreadyIncluded = true
				if len(g.Applys) == 0 {
					m.srcGroups[i] = m.srcGroups[len(m.srcGroups)-1]
					m.srcGroups = m.srcGroups[:len(m.srcGroups)-1]
				}
			}
		}
		if !alreadyIncluded && len(msg.Group.Applys) > 0 {
			m.srcGroups = append(m.srcGroups, msg.Group)
		}
		m.groups = make([]tea.Model, len(m.srcGroups))
		for i, g := range m.srcGroups {
			m.groups[i] = apply.New(g, i, m.width-2)
			cmd = m.groups[i].Init()
		}
		m.redrawContent()
	case ActiveMsg:
		m.active = msg.Active
		m.updateActive()
	case apply.ChangeMsg:
		m.groups[msg.Gidx], cmd = m.groups[msg.Gidx].Update(msg)
		m.redrawContent()
	case apply.SuccessMsg:
		m.groups[msg.Gidx], cmd = m.groups[msg.Gidx].Update(msg)
		m.redrawContent()
	case apply.FailMsg:
		m.groups[msg.Gidx], cmd = m.groups[msg.Gidx].Update(msg)
		m.redrawContent()
	}
	m.updateView()
	if m.active {
		if cmd == nil {
			cmd = m.cursorPosition
		} else {
			cmd = tea.Batch(cmd, m.cursorPosition)
		}

	}
	return m, cmd
}

func (m *model) cursorUp() {
	m.selectedGroup--
	if m.selectedGroup < 0 {
		m.selectedGroup = 0
	}
}
func (m *model) cursorDown() {
	m.selectedGroup++
	if m.selectedGroup >= len(m.groups) {
		m.selectedGroup = len(m.groups) - 1
	}
}

func (m *model) cursorPosition() tea.Msg {
	return separator.CursorMsg{
		PointingApply: true,
		Position:      m.selectedGroupFirstLine,
	}
}
