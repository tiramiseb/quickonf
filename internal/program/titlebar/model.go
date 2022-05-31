package titlebar

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/button"
	"github.com/tiramiseb/quickonf/internal/program/messages"
	"github.com/tiramiseb/quickonf/internal/program/styles"
)

var style = lipgloss.NewStyle().
	Background(styles.FgColor).
	Foreground(styles.BgColor).
	Bold(true)

type Model struct {
	title       string
	middleSpace int

	apply      *button.Button
	applyStart int
	applyEnd   int
	showApply  bool

	recheck      *button.Button
	recheckStart int
	recheckEnd   int
	showRecheck  bool

	filter      *button.Toggle
	filterStart int
	filterEnd   int
	showFilter  bool

	details      *button.Toggle
	detailsStart int
	detailsEnd   int
	showDetails  bool

	help      *button.Button
	helpStart int
	helpEnd   int
	showHelp  bool

	quit      *button.Button
	quitStart int
	quitEnd   int
	showQuit  bool

	helpBack      *button.Button
	helpBackStart int
	helpBackEnd   int

	HelpView func() string
}

func New() *Model {
	return &Model{
		title:   " Quickonf ",
		apply:   button.NewButton("Apply", 0, apply),
		recheck: button.NewButton("Recheck", 0, recheck),
		filter:  button.NewToggle("Filter checks", 0, "filter", true),
		details: button.NewToggle("More details", 5, "details", false),
		help:    button.NewButton("Help", 0, enableHelp),
		quit:    button.NewButton("Quit", 0, tea.Quit),

		helpBack: button.NewButton("Back (esc)", -2, disableHelp),

		HelpView: func() string { return "" },
	}
}

func enableHelp() tea.Msg {
	return messages.Toggle{Name: "help", Action: messages.ToggleActionEnable}
}

func disableHelp() tea.Msg {
	return messages.Toggle{Name: "help", Action: messages.ToggleActionDisable}
}

func recheck() tea.Msg {
	return messages.Recheck{}
}

func apply() tea.Msg {
	return messages.Apply{}
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
			switch {
			case msg.X >= m.applyStart && msg.X <= m.applyEnd:
				cmd = m.apply.Click
			case msg.X >= m.recheckStart && msg.X <= m.recheckEnd:
				cmd = m.recheck.Click
			case msg.X >= m.filterStart && msg.X <= m.filterEnd:
				cmd = m.filter.Click
			case msg.X >= m.detailsStart && msg.X <= m.detailsEnd:
				cmd = m.details.Click
			case msg.X >= m.helpStart && msg.X <= m.helpEnd:
				cmd = m.help.Click
			case msg.X >= m.quitStart && msg.X <= m.quitEnd:
				cmd = m.quit.Click
			}
		}
	case messages.ToggleStatus:
		switch msg.Name {
		case "filter":
			m.filter = m.filter.ChangeStatus(msg.Status)
		case "details":
			m.details = m.details.ChangeStatus(msg.Status)
		}

	}
	return m, cmd
}

func (m *Model) UpdateInHelp(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type == tea.MouseRelease {
			switch {
			case msg.X >= m.helpBackStart && msg.X <= m.helpBackEnd:
				cmd = m.helpBack.Click
			}
		}
	case messages.ToggleStatus:
		switch msg.Name {
		case "filter":
			m.filter = m.filter.ChangeStatus(msg.Status)
		case "details":
			m.details = m.details.ChangeStatus(msg.Status)
		}

	}
	return m, cmd
}
