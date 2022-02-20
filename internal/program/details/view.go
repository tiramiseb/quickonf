package details

func (m *Model) RedrawView() *Model {
	m.completeView = "details"
	// TODO Draw
	return m
}

func (m *Model) View() string {
	return m.style.Render(m.completeView)
}
