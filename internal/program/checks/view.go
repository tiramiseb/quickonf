package checks

import (
	"fmt"
	"strings"
)

func (m *Model) prepareView(width int) {
	lines := make([]string, len(m.groups))
	for i, g := range m.groups {
		line := fmt.Sprintf("%d: %s", g.Status(), g.Name)
		if len(line) > width {
			line = line[:width-1]
		}
		lines[i] = line

	}
	m.completeView = lines
}

func (m *Model) RedrawGroup(i int) *Model {
	g := m.groups[i]
	m.completeView[i] = fmt.Sprintf("%d: %s", g.Status(), g.Name)
	return m
}

func (m *Model) View() string {
	return m.style.Render(strings.Join(m.completeView, "\n"))
}
