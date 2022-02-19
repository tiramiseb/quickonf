package titlebar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/button"
	"github.com/tiramiseb/quickonf/internal/program/messages"
)

const (
	bgStyle = lipgloss.Color("#555555")
	fgStyle = lipgloss.Color("#ffffff")
)

var style = lipgloss.NewStyle().
	Background(bgStyle).
	Foreground(fgStyle).
	Bold(true)

type Model struct {
	// toggle      tea.Model
	// toggleStart int
	// toggleEnd   int

	filter      *button.Toggle
	filterStart int
	filterEnd   int

	help      *button.Toggle
	helpStart int
	helpEnd   int

	quit      *button.Button
	quitStart int
	quitEnd   int

	isInHelp bool

	view     func() string
	helpView func() string
}

func New() *Model {
	return &Model{
		// toggle:   button.New("Toggle", 0, messages.Toggle),
		filter:   button.NewToggle("Filter", 0, messages.Filter(true), messages.Filter(false), true),
		help:     button.NewToggle("Help", 0, messages.Help(true), messages.Help(false), false),
		quit:     button.NewButton("Quit", 0, tea.Quit),
		view:     func() string { return "" },
		helpView: func() string { return "" },
	}
}

// Resize resizes the titlebar
func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	m.draw(size.Width)
	return m
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch {
		case msg.X >= m.filterStart && msg.X <= m.filterEnd:
			m.filter, cmd = m.filter.Click()
		case msg.X >= m.helpStart && msg.X <= m.helpEnd:
			m.help, cmd = m.help.Click()
		case msg.X >= m.quitStart && msg.X <= m.quitEnd:
			m.quit, cmd = m.quit.Click()
		}
		// case messages.FilterMsg:
		// 	m.filter, cmd = m.filter.Update(togglebutton.Toggle{On: msg.On})
		// case messages.HelpMsg:
		// 	m.isInHelp = msg.On
		// 	m.help, cmd = m.help.Update(togglebutton.Toggle{On: msg.On})
	}
	return m, cmd
}

func (m *Model) HelpActive(show bool) *Model {
	m.help = m.help.FromExternal(show)
	return m
}
