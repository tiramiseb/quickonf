package titlebar

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/internal/program/button"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

type model struct {
	quit      tea.Model
	quitStart int
	quitEnd   int
	view      func() string
}

func New() *model {
	quit := button.New("Quit", 0, tea.Quit)
	return &model{
		quit: quit,
		view: func() string { return "" },
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "Q", "esc":
			cmd = tea.Quit
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseUnknown {
			m.quit, _ = m.quit.Update(msg)
			break
		}
		unknown := tea.MouseMsg{Type: tea.MouseUnknown}
		switch {
		case msg.X >= m.quitStart && msg.X <= m.quitEnd:
			m.quit, cmd = m.quit.Update(msg)
		default:
			m.quit, _ = m.quit.Update(unknown)
		}
	}
	return m, cmd
}

func (m *model) makeView(width int) {
	title := " Quickonf "
	quitWidth := lipgloss.Width(m.quit.View())
	spaceWidth := width - 10 - quitWidth - 1
	switch {
	case width < 10:
		view := style.Title.Render(title[:width])
		m.view = func() string {
			return view
		}
		m.quitStart = -1
		m.quitEnd = -1
	case spaceWidth < 0:
		view := style.Title.Render(title + strings.Repeat(" ", width-10))
		m.view = func() string {
			return view
		}
		m.quitStart = -1
		m.quitEnd = -1
	default:
		leftPart := style.Title.Render(title + strings.Repeat(" ", spaceWidth))
		space := style.Title.Render(" ")
		m.view = func() string {
			return leftPart + m.quit.View() + space
		}
		// Point to the last char of the quit button
		i := width - 2

		m.quitEnd = i
		i -= lipgloss.Width(m.quit.View()) - 1
		m.quitStart = i

		i--
		i--
	}
}

func (m *model) View() string {
	return m.view()
}
