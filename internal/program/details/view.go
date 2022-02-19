package details

func (m *Model) prepareView() {
	m.completeView = m.style.Render("details")
	// TODO Draw
}

func (m *Model) View() string {
	return m.completeView
}
