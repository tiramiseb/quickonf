package header

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/style"
)

const Height = 3

type Model struct {
	width int
	title string

	View string
}

func New(filtered bool, nb int) *Model {
	var title string
	if filtered {
		if nb > 1 {
			title = fmt.Sprintf("Quickonf: Applying %d groups", nb)
		} else {
			title = "Quickonf: Applying 1 group"

		}
	} else {
		title = "Quickonf: Applying all groups"
	}
	return &Model{title: title}
}

func (m *Model) Update(msg tea.Msg) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.update()
	}
}

func (m *Model) update() {
	title := m.title
	if len(m.title) > m.width {
		title = title[:m.width]
	}
	m.View = style.Header.
		Width(m.width).
		Render(title)
}
