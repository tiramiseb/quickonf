package program

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/internal/program/footer"
	"github.com/tiramiseb/quickonf/internal/program/groupslist"
	"github.com/tiramiseb/quickonf/internal/program/header"
	"github.com/tiramiseb/quickonf/internal/program/style"
	"github.com/tiramiseb/quickonf/internal/program/toppanel"
	"github.com/tiramiseb/quickonf/internal/state"
)

type model struct {
	state *state.State

	header     *header.Model
	footer     *footer.Model
	toppanel   *toppanel.Model
	groupslist *groupslist.Model
}

const verticalMargins = header.Height + footer.Height

func newModel(st *state.State) *model {
	nb := len(st.Groups)
	return &model{
		state: st,

		header:     header.New(st.Filtered, nb),
		footer:     footer.New(nb),
		toppanel:   toppanel.New(80),
		groupslist: groupslist.New(80, 24, verticalMargins, st),
	}
}

func (m *model) Init() tea.Cmd {
	tea.LogToFile("/tmp/tmplog", "")
	return m.groupslist.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	default:
		m.header.Update(msg)
		m.footer.Update(msg)
	}
	return m, tea.Batch(
		m.toppanel.Update(msg),
		m.groupslist.Update(msg),
	)
}

func (m *model) View() string {
	return style.Main.Render(
		fmt.Sprintf(
			"%s\n%s%s\n%s",
			m.header.View,
			m.toppanel.View,
			m.groupslist.View(m.toppanel.Size),
			m.footer.View,
		),
	)
}
