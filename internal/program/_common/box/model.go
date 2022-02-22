package box

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/specific/separator"
)

type UpdateElementsMsg struct {
	Elements []tea.Model
}

type ForceRedrawMsg struct{}

type ElementSelectedMsg struct {
	Selected bool
}

type elementLine struct {
	idx         int
	elementline int
}

type model struct {
	title             string
	msgIfEmpty        string
	elements          []tea.Model
	cursorPointsRight bool // Does separator cursor point right (if false, it points left) when this box is active?

	width         int
	boxHeight     int
	active        bool
	subtitleStyle lipgloss.Style
	boxStyle      lipgloss.Style

	allElementsView          []string      // All elements! And then take a window of this for view
	allLineToElement         []elementLine // map of line number to displayed elements, for passing clicks to the correct element
	selectedElement          int           // index of the selected element in the elements list
	selectedElementFirstLine int           // line on screen of the first line of the selected element (for cursor position)
	view                     string        // The resulting view itself
	viewLineToElement        []elementLine
}

func New(title, msgIfEmpty string, elements []tea.Model, cursorPointsRight, active bool) tea.Model {
	return &model{
		title:             title,
		msgIfEmpty:        msgIfEmpty,
		elements:          elements,
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
	case UpdateElementsMsg:
		m.elements = msg.Elements
		m.redrawContent()
	case tea.WindowSizeMsg:
		cmd = m.windowSize(msg)
		m.redrawContent()
	case tea.KeyMsg:
		if len(m.elements) == 0 {
			m.redrawContent()
			return m, nil
		}
		switch msg.String() {
		case "up":
			cmd = m.cursorUp()
			m.redrawContent()
		case "down":
			cmd = m.cursorDown()
			m.redrawContent()
		case "home":
			cmd = m.cursorFirst()
			m.redrawContent()
		case "end":
			cmd = m.cursorLast()
			m.redrawContent()
		case "pgup":
			cmd = m.cursorWindowUp()
			m.redrawContent()
		case "pgdown":
			cmd = m.cursorWindowDown()
			m.redrawContent()
		default:
			m.elements[m.selectedElement], cmd = m.elements[m.selectedElement].Update(msg)
			m.redrawContent()
		}
	case tea.MouseMsg:
		// MouseUnknown means the user clicked somewhere else, inform everyone
		cmds := make([]tea.Cmd, len(m.elements))
		if msg.Type == tea.MouseUnknown {
			for i, g := range m.elements {
				m.elements[i], cmds[i] = g.Update(msg)
			}
			cmd = tea.Batch(cmds...)
			m.redrawContent()
			break
		}

		// These situations mean the event is not over an element
		if msg.X <= 0 || msg.X >= m.width-1 || msg.Y <= 1 || msg.Y >= m.boxHeight+2 {
			unknown := tea.MouseMsg{
				Type: tea.MouseUnknown,
			}
			for i, g := range m.elements {
				m.elements[i], cmds[i] = g.Update(unknown)
			}
			cmd = tea.Batch(cmds...)
			m.redrawContent()
			break
		}

		// Mouse wheel movement
		if msg.Type == tea.MouseWheelUp {
			cmd = m.cursorUp()
			m.redrawContent()
			break
		}
		if msg.Type == tea.MouseWheelDown {
			cmd = m.cursorDown()
			m.redrawContent()
			break
		}

		// Forward mouse message to the element under the cursor
		unknown := tea.MouseMsg{
			Type: tea.MouseUnknown,
		}
		msg.X--
		line := msg.Y - 2
		clickedElement := -1
		if line < len(m.viewLineToElement) {
			elementline := m.viewLineToElement[line]
			if msg.Type == tea.MouseRelease && elementline.idx != m.selectedElement {
				var cmd1, cmd2 tea.Cmd
				m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
				m.selectedElement = elementline.idx
				m.elements[elementline.idx], cmd2 = m.elements[elementline.idx].Update(ElementSelectedMsg{true})
				cmds[elementline.idx] = tea.Batch(cmd1, cmd2)
			} else {
				msg.Y = elementline.elementline
				m.elements[elementline.idx], cmds[elementline.idx] = m.elements[elementline.idx].Update(msg)
				clickedElement = elementline.idx
			}
		}

		// And provide unknown to all other elements
		for i, g := range m.elements {
			if i == clickedElement {
				continue
			}
			m.elements[i], cmds[i] = g.Update(unknown)
		}
		cmd = tea.Batch(cmds...)
		m.redrawContent()
	case separator.ActiveMsg:
		isActive := msg.IsRightActive == m.cursorPointsRight
		m.active = isActive
		m.elements[m.selectedElement], cmd = m.elements[m.selectedElement].Update(ElementSelectedMsg{isActive})
		m.updateActive()
		m.redrawContent()
	case ForceRedrawMsg:
		m.redrawContent()
	default:
		m.elements[m.selectedElement], cmd = m.elements[m.selectedElement].Update(msg)
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