package titlebar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/button"
	"github.com/tiramiseb/quickonf/internal/program/global"
	"github.com/tiramiseb/quickonf/internal/program/global/toggles"
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

	help      *button.Button
	helpStart int
	helpEnd   int

	quit      *button.Button
	quitStart int
	quitEnd   int

	helpBack      *button.Button
	helpBackStart int
	helpBackEnd   int

	view     func() string
	helpView func() string
}

func New() *Model {
	return &Model{
		filter:  button.NewToggle("Filter checks", 0, "filter"),
		details: button.NewToggle("More details", 5, "details"),
		help:    button.NewButton("Help", 0, toggleHelp),
		quit:    button.NewButton("Quit", 0, tea.Quit),

		helpBack: button.NewButton("Back (esc)", -2, toggleHelp),

		view:     func() string { return "" },
		helpView: func() string { return "" },
	}
}

func toggleHelp() tea.Msg {
	toggles.Toggle("help")
	return nil
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
		if msg.Type == tea.MouseRelease {
			if toggles.Get("help") {
				switch {
				case msg.X >= m.helpBackStart && msg.X <= m.helpBackEnd:
					cmd = m.helpBack.Click()
				}
			} else {
				switch {
				case msg.X >= m.filterStart && msg.X <= m.filterEnd:
					cmd = m.filter.Click()
				case msg.X >= m.detailsStart && msg.X <= m.detailsEnd:
					cmd = m.details.Click()
				case msg.X >= m.helpStart && msg.X <= m.helpEnd:
					cmd = m.help.Click()
				case msg.X >= m.quitStart && msg.X <= m.quitEnd:
					cmd = m.quit.Click()
				}
			}
		}
	}
	return m, cmd
}
