package details

func (m *Model) ChangeView(idx int) *Model {
	m.displayedGroup = idx
	return m
}

func (m *Model) View() string {
	if m.displayedGroup < 0 {
		m.displayedGroup = 0
	} else if m.displayedGroup >= len(m.groups) {
		m.displayedGroup = len(m.groups) - 1
	}
	return m.style.Render("There will be details for " + m.groups[m.displayedGroup].Name)
}
