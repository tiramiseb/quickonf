package reallyreloadconfig

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tiramiseb/quickonf/program/button"
	"github.com/tiramiseb/quickonf/program/messages"
)

// Dialog is 6 lines high.
// Dialog is 57 characters wide.
// There are 19 spaces before the buttons.
// Buttons are on the 5th line.
// There are 10 spaces between the buttons.
const dialogTemplate = `
Are you sure you really want to reload the configuration?


%s          %s
`

type Model struct {
	yes      *button.Button
	yesStart int
	yesEnd   int

	no      *button.Button
	noStart int
	noEnd   int

	buttonsY int

	dialog string
	view   string
}

func New() *Model {
	yes := button.NewButton("Yes", 0, yes)
	no := button.NewButton("No", 0, no)
	dialog := fmt.Sprintf(dialogTemplate, yes.View, no.View)
	return &Model{
		yes:    yes,
		no:     no,
		dialog: dialog,
	}
}

func yes() tea.Msg {
	return messages.ConfirmReloadConfig{}
}

func no() tea.Msg {
	return messages.ReloadConfig{}
}

// Resize resizes the "really apply?" view
func (m *Model) Resize(size tea.WindowSizeMsg) *Model {
	if size.Width < 50 {
		size.Width = 50
	}
	if size.Height < 7 {
		size.Height = 7
	}
	m.buttonsY = (size.Height-6)/2 + 4
	m.yesStart = (size.Width-57)/2 + 19
	m.yesEnd = m.yesStart + m.yes.Width - 1
	m.noStart = m.yesEnd + 11
	m.noEnd = m.noStart + m.no.Width - 1

	m.view = lipgloss.Place(size.Width, size.Height, lipgloss.Center, lipgloss.Center, m.dialog)
	return m
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type != tea.MouseRelease || msg.Y != m.buttonsY {
			break
		}
		switch {
		case msg.X >= m.yesStart && msg.X <= m.yesEnd:
			cmd = m.yes.Click
		case msg.X >= m.noStart && msg.X <= m.noEnd:
			cmd = m.no.Click
		}

	}
	return m, cmd
}

func (m *Model) View() string {
	return m.view
}
