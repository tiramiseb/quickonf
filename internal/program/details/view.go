package details

func (m *Model) prepareView() {
	m.completeView = "details"
	// TODO Draw
}

func (m *Model) View() string {
	return m.style.Render(m.completeView)
}
