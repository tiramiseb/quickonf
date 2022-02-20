package titlebar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/button"
	"github.com/tiramiseb/quickonf/internal/program/global"
	"github.com/tiramiseb/quickonf/internal/program/messages"
)

var style = lipgloss.NewStyle().
	Background(global.FgColor).
	Foreground(global.BgColor).
	Bold(true)

type Model struct {
	filter      *button.Toggle
	filterStart int
	filterEnd   int

	details      *button.Toggle
	detailsStart int
	detailsEnd   int

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
		filter:   button.NewToggle("Filter checks", 0, messages.Filter, "filter"),
		details:  button.NewToggle("More details", 5, messages.Details, "details"),
		help:     button.NewToggle("Help", 0, messages.Help, "help"),
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
		case msg.X >= m.detailsStart && msg.X <= m.detailsEnd:
			m.details, cmd = m.details.Click()
		case msg.X >= m.helpStart && msg.X <= m.helpEnd:
			m.help, cmd = m.help.Click()
		case msg.X >= m.quitStart && msg.X <= m.quitEnd:
			m.quit, cmd = m.quit.Click()
		}
	}
	return m, cmd
}
