package box

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/common/group"
	"github.com/tiramiseb/quickonf/internal/program/specific/separator"
)

type UpdateGroupsMsg struct {
	Groups []tea.Model
}

type groupLine struct {
	gidx      int
	groupline int
}

type model struct {
	title             string
	msgIfEmpty        string
	groups            []tea.Model
	cursorPointsRight bool // Does separator cursor point right (if false, it points left) when this box is active?

	width         int
	boxHeight     int
	active        bool
	subtitleStyle lipgloss.Style
	boxStyle      lipgloss.Style

	allGroupsView          []string    // All groups! And then take a window of this for view
	allLineToGroup         []groupLine // map of line number to displayed group, for passing clicks to the correct group
	selectedGroup          int         // index of the selected group in the groups list
	selectedGroupFirstLine int         // line on screen of the first line of the selected group (for cursor position)
	view                   string      // The resulting view itself
	viewLineToGroup        []groupLine
}

func New(title, msgIfEmpty string, groups []tea.Model, cursorPointsRight, active bool) tea.Model {
	return &model{
		title:             title,
		msgIfEmpty:        msgIfEmpty,
		groups:            groups,
		cursorPointsRight: cursorPointsRight,
		width:             2,
		boxHeight:         1,
		active:            active,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case UpdateGroupsMsg:
		m.groups = msg.Groups
		m.redrawContent()
	case tea.WindowSizeMsg:
		cmd = m.windowSize(msg)
		m.redrawContent()
	case tea.KeyMsg:
		if len(m.groups) == 0 {
			m.redrawContent()
			return m, nil
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
			groupline := m.viewLineToGroup[line]
			msg.Y = groupline.groupline
			m.groups[groupline.gidx], cmds[groupline.gidx] = m.groups[groupline.gidx].Update(msg)
			clickedGroup = groupline.gidx
		}

		// And provide unknown to all other groups
		for i, g := range m.groups {
			if i == clickedGroup {
				continue
			}
			m.groups[i], cmds[i] = g.Update(unknown)
		}
		cmd = tea.Batch(cmds...)
		m.redrawContent()
	case separator.ActiveMsg:
		m.active = msg.IsRightActive == m.cursorPointsRight
		m.updateActive()
	case group.Msg:
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
