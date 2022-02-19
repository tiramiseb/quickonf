package help

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/commands"
	"github.com/tiramiseb/quickonf/internal/program/common/messages"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
)

type section int

const (
	sectionIntro section = iota
	sectionLang
	sectionCommands
	sectionUI
)

const (
	helpOnHelp = "To close help, press Esc "
	helpLen    = len(helpOnHelp)
)

type model struct {
	introStart    int
	introEnd      int
	langStart     int
	langEnd       int
	commandsStart int
	commandsEnd   int
	uiStart       int
	uiEnd         int

	width        int
	buttonsWidth int

	activeIntro    string
	activeLang     string
	activeCommands string
	activeUI       string

	helpOnHelp string

	shallRender bool

	active     section
	boxStyle   lipgloss.Style
	viewport   viewport.Model
	mdRenderer *glamour.TermRenderer

	showFilter      bool
	filterMaxLength int
	filter          string
	commands        []commands.Command
	contents        map[section]string

	filteredCommandsDoc map[string]string // Map of filter to rendered doc
}

func New() *model {
	r, _ := glamour.NewTermRenderer()
	m := &model{
		width:               8,
		shallRender:         true,
		viewport:            viewport.Model{Width: 8, Height: 1},
		mdRenderer:          r,
		commands:            commands.GetAll(),
		contents:            map[section]string{},
		filteredCommandsDoc: map[string]string{},
	}
	m.viewport.SetContent(m.contents[m.active])
	return m
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.boxStyle = style.ActiveBox.Copy().Width(msg.Width - 2).Height(msg.Height - 4)
		m.viewport.Height = msg.Height - 4
		m.viewport.Width = msg.Width - 2
		m.updateNav(msg.Width)
		var helpFullWidth string
		if helpLen >= msg.Width {
			helpFullWidth = helpOnHelp[helpLen-msg.Width:]
		} else {
			helpFullWidth = strings.Repeat(" ", m.width-helpLen) + helpOnHelp
		}
		m.helpOnHelp = style.BoxContent.Render(helpFullWidth)
		m.mdRenderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(msg.Width-2),
		)
		m.shallRender = true
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "left":
			m.active--
			if m.active < sectionIntro {
				m.active = sectionIntro
			}
			m.viewport.SetContent(m.contents[m.active])
		case "right":
			m.active++
			if m.active > sectionUI {
				m.active = sectionUI
			}
			m.viewport.SetContent(m.contents[m.active])
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
			"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
			"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
			"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
			".":
			if m.active == sectionCommands && len(m.filter) < m.filterMaxLength {
				m.filter += strings.ToLower(key)
				m.contents[sectionCommands] = m.commandsDoc()
				m.viewport.SetContent(m.contents[sectionCommands])
			}
		case "backspace":
			if m.active == sectionCommands && len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.contents[sectionCommands] = m.commandsDoc()
				m.viewport.SetContent(m.contents[sectionCommands])
			}
		default:
			m.viewport, cmd = m.viewport.Update(msg)
		}
	case tea.MouseMsg:
		switch {
		case msg.Y != 0 || msg.Type != tea.MouseRelease:
			m.viewport, cmd = m.viewport.Update(msg)
		case msg.X >= m.introStart && msg.X <= m.introEnd:
			m.active = sectionIntro
			m.viewport.SetContent(m.contents[m.active])
		case msg.X >= m.langStart && msg.X <= m.langEnd:
			m.active = sectionLang
			m.viewport.SetContent(m.contents[m.active])
		case msg.X >= m.commandsStart && msg.X <= m.commandsEnd:
			m.active = sectionCommands
			m.viewport.SetContent(m.contents[m.active])
		case msg.X >= m.uiStart && msg.X <= m.uiEnd:
			m.active = sectionUI
			m.viewport.SetContent(m.contents[m.active])
		}
	case messages.HelpMsg:
		if msg.On && m.shallRender {
			m.shallRender = false
			m.reRenderContents()
		}
	}
	return m, cmd
}

func (m *model) View() string {
	var toolbarLeft string
	switch m.active {
	case sectionIntro:
		toolbarLeft = m.activeIntro
	case sectionLang:
		toolbarLeft = m.activeLang
	case sectionCommands:
		toolbarLeft = m.activeCommands
	case sectionUI:
		toolbarLeft = m.activeUI
	}

	var (
		toolbarRight string
		filterWidth  int
	)
	if m.showFilter && m.active == sectionCommands {
		if m.filter == "" {
			toolbarRight = style.BoxContent.Render("Filter: ")
			filterWidth = 8
		} else {
			toolbarRight = style.ClickedButton.Render("Filter: "+m.filter) + style.BoxContent.Render(" ")
			filterWidth = 9 + len(m.filter)
		}
	}

	var spaces string
	nbSpaces := m.width - m.buttonsWidth - filterWidth
	if nbSpaces > 0 {
		spaces = strings.Repeat(" ", nbSpaces)
	}

	return toolbarLeft + style.BoxContent.Render(spaces) + toolbarRight + "\n" + m.boxStyle.Render(m.viewport.View()) + "\n" + m.helpOnHelp
}
