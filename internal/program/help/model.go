package help

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/tiramiseb/quickonf/internal/commands"
)

const filterMaxLength = 20

type Model struct {
	viewport viewport.Model

	height int
	width  int

	introTitle             string
	introTitleWithFocus    string
	introStart             int
	introEnd               int
	languageTitle          string
	languageTitleWithFocus string
	languageStart          int
	languageEnd            int
	commandsTitle          string
	commandsTitleWithFocus string
	commandsStart          int
	commandsEnd            int
	uiTitle                string
	uiTitleWithFocus       string
	uiStart                int
	uiEnd                  int
	subtitleSeparator      string

	activeSection       int
	commandFilter       string
	filteredCommandsDoc map[string]string
	commands            []commands.Command
}

func New() *Model {
	v := viewport.Model{Width: 1, Height: 1}
	v.Style = lipgloss.NewStyle().MarginLeft(2).MarginRight(2)
	return &Model{
		viewport:            v,
		activeSection:       0,
		filteredCommandsDoc: map[string]string{},
		commands:            commands.GetAll(),
	}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "left":
			if m.activeSection > 0 {
				m.activeSection--
				m.setContent()
			}
		case "right":
			if m.activeSection < 3 {
				m.activeSection++
				m.setContent()
			}
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
			"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			".":
			if m.activeSection == 2 && len(m.commandFilter) < filterMaxLength {
				m.commandFilter += strings.ToLower(key)
				m.setContent()
			}
		case "backspace":
			if m.activeSection == 2 && len(m.commandFilter) > 0 {
				m.commandFilter = m.commandFilter[:len(m.commandFilter)-1]
				m.setContent()
			}
		default:
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}
	return m, cmd
}

func (m *Model) setContent() {
	m.viewport.Width = m.width
	var text string
	if lipgloss.HasDarkBackground() {
		switch m.activeSection {
		case 0:
			text = introDark
			m.viewport.Height = m.height - 2
		case 1:
			text = languageDark
			m.viewport.Height = m.height - 2
		case 2:
			m.viewport.Height = m.height - 4
			text = m.commandsDoc(true)
		case 3:
			text = uiDark
			m.viewport.Height = m.height - 2
		}
	} else {
		switch m.activeSection {
		case 0:
			text = introLight
			m.viewport.Height = m.height - 2
		case 1:
			text = languageLight
			m.viewport.Height = m.height - 2
		case 2:
			m.viewport.Height = m.height - 4
			text = m.commandsDoc(false)
		case 3:
			text = uiLight
			m.viewport.Height = m.height - 2
		}
	}
	m.viewport.SetContent(
		wordwrap.String(text, m.width-4),
	)
}
