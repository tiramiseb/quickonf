package titlebar

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tiramiseb/quickonf/internal/program/common/button"
	"github.com/tiramiseb/quickonf/internal/program/common/messages"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
	"github.com/tiramiseb/quickonf/internal/program/common/togglebutton"
)

type model struct {
	toggle      tea.Model
	toggleStart int
	toggleEnd   int

	filter      tea.Model
	filterStart int
	filterEnd   int

	help      tea.Model
	helpStart int
	helpEnd   int

	quit      tea.Model
	quitStart int
	quitEnd   int

	isInHelp bool
	view     func() string
	helpView func() string
}

func New() *model {
	return &model{
		toggle:   button.New("Toggle", 0, messages.Toggle),
		filter:   togglebutton.New("Filter", 0, messages.Filter(true), messages.Filter(false)),
		help:     togglebutton.New("Help", 0, messages.Help(true), messages.Help(false)),
		quit:     button.New("Quit", 0, tea.Quit),
		view:     func() string { return "" },
		helpView: func() string { return "" },
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.makeView(msg.Width)
	case tea.MouseMsg:
		if msg.Type == tea.MouseUnknown {
			m.quit, _ = m.quit.Update(msg)
			break
		}
		unknown := tea.MouseMsg{Type: tea.MouseUnknown}
		switch {
		case msg.X >= m.toggleStart && msg.X <= m.toggleEnd:
			m.toggle, cmd = m.toggle.Update(msg)
			m.filter, _ = m.filter.Update(unknown)
			m.help, _ = m.help.Update(unknown)
			m.quit, _ = m.quit.Update(unknown)
		case msg.X >= m.filterStart && msg.X <= m.filterEnd:
			m.toggle, _ = m.toggle.Update(unknown)
			m.filter, cmd = m.filter.Update(msg)
			m.help, _ = m.help.Update(unknown)
			m.quit, _ = m.quit.Update(unknown)
		case msg.X >= m.helpStart && msg.X <= m.helpEnd:
			m.toggle, _ = m.toggle.Update(unknown)
			m.filter, _ = m.filter.Update(unknown)
			m.help, cmd = m.help.Update(msg)
			m.quit, _ = m.quit.Update(unknown)
		case msg.X >= m.quitStart && msg.X <= m.quitEnd:
			m.toggle, _ = m.toggle.Update(unknown)
			m.filter, _ = m.filter.Update(unknown)
			m.help, _ = m.help.Update(unknown)
			m.quit, cmd = m.quit.Update(msg)
		default:
			m.toggle, _ = m.toggle.Update(unknown)
			m.filter, _ = m.filter.Update(unknown)
			m.help, _ = m.help.Update(unknown)
			m.quit, _ = m.quit.Update(unknown)
		}
	case messages.FilterMsg:
		m.filter, cmd = m.filter.Update(togglebutton.Toggle{On: msg.On})
	case messages.HelpMsg:
		m.isInHelp = msg.On
		m.help, cmd = m.help.Update(togglebutton.Toggle{On: msg.On})
	}
	return m, cmd
}

func (m *model) makeView(width int) {
	title := " Quickonf "
	titleWidth := len(title)
	m.quitStart = -1
	m.quitEnd = -1
	m.helpStart = -1
	m.helpEnd = -1
	m.filterStart = -1
	m.filterEnd = -1
	m.toggleStart = -1
	m.toggleEnd = -1
	availableWidth := width - titleWidth

	// No place for the quit button, only include the title
	buttonWidth := lipgloss.Width(m.quit.View())
	if availableWidth <= buttonWidth {
		var view string

		if availableWidth <= 0 {
			view = style.Title.Render(title[:width])
		} else {
			view = style.Title.Render(title + strings.Repeat(" ", availableWidth))
		}
		m.view = func() string {
			return view
		}
		return
	}
	space := style.Title.Render(" ")
	availableWidth = availableWidth - buttonWidth - 1
	m.quitEnd = width - 2
	m.quitStart = m.quitEnd - buttonWidth + 1

	// No place for the help button, include title & quit
	buttonWidth = lipgloss.Width(m.help.View())
	if availableWidth <= buttonWidth {
		leftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
		m.view = func() string {
			return leftAndSpace + m.quit.View() + space
		}
		m.helpView = m.view
		return
	}
	availableWidth = availableWidth - buttonWidth - 1
	m.helpEnd = m.quitStart - 2
	m.helpStart = m.helpEnd - buttonWidth + 1

	helpLeftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
	m.helpView = func() string {
		return helpLeftAndSpace + m.help.View() + space + m.quit.View() + space
	}

	// No place for the filter button, include title & help & quit
	buttonWidth = lipgloss.Width(m.filter.View())
	if availableWidth <= buttonWidth {
		leftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
		m.view = func() string {
			return leftAndSpace + m.help.View() + space + m.quit.View() + space
		}
		return
	}

	availableWidth = availableWidth - buttonWidth - 1
	m.filterEnd = m.helpStart - 2
	m.filterStart = m.filterEnd - buttonWidth + 1

	// No place for the toggle button, include title & filter & help & quit
	buttonWidth = lipgloss.Width(m.toggle.View())
	if availableWidth <= buttonWidth {
		leftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
		m.view = func() string {
			return leftAndSpace + m.filter.View() + space + m.help.View() + space + m.quit.View() + space
		}
		return
	}
	availableWidth = availableWidth - buttonWidth - 1
	m.toggleEnd = m.filterStart - 2
	m.toggleStart = m.toggleEnd - buttonWidth + 1

	leftAndSpace := style.Title.Render(title + strings.Repeat(" ", availableWidth))
	m.view = func() string {
		return leftAndSpace + m.toggle.View() + space + m.filter.View() + space + m.help.View() + space + m.quit.View() + space
	}
}

func (m *model) View() string {
	if m.isInHelp {
		return m.helpView()
	}
	return m.view()
}
