package program

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/tiramiseb/quickonf/program/footer"
	"github.com/tiramiseb/quickonf/program/groupslist"
	"github.com/tiramiseb/quickonf/program/header"
	"github.com/tiramiseb/quickonf/state"
)

type model struct {
	state *state.State

	header     *header.Model
	footer     *footer.Model
	groupslist *groupslist.Model
}

const verticalMargins = header.Height + footer.Height

func newModel(st *state.State) *model {
	nb := len(st.Groups)
	return &model{
		state: st,

		header:     header.New(st.Filtered, nb),
		footer:     footer.New(nb),
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
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}
	default:
		m.header.Update(msg)
		m.footer.Update(msg)
	}
	return m, m.groupslist.Update(msg)
}

func (m *model) View() string {
	return fmt.Sprintf("%s\n%s\n%s", m.header.View, m.groupslist.View(), m.footer.View)
}
