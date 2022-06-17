package groups

const scrollDelta = 3 // Same value as default value for viewport's mousewheeldelta

func (m *Model) up() {
	m.selectedGroup = m.selectedGroup.Previous(1, m.showSuccessful)
}

func (m *Model) down() {
	m.selectedGroup = m.selectedGroup.Next(1, m.showSuccessful)
}

func (m *Model) pgup() {
	m.selectedGroup = m.selectedGroup.Previous(m.height/2, m.showSuccessful)
}

func (m *Model) pgdown() {
	m.selectedGroup = m.selectedGroup.Next(m.height/2, m.showSuccessful)
}

func (m *Model) home() {
	m.selectedGroup = m.selectedGroup.Previous(m.groups.Count(), m.showSuccessful)
}

func (m *Model) end() {
	m.selectedGroup = m.selectedGroup.Next(m.groups.Count(), m.showSuccessful)
}

func (m *Model) selectLine(lineIdx int) {
	m.selectedGroup = m.firstDisplayedGroup.Next(lineIdx, m.showSuccessful)
}

func (m *Model) scrollUp() {
	m.selectedGroup = m.selectedGroup.Previous(scrollDelta, m.showSuccessful)
}

func (m *Model) scrollDown() {
	m.selectedGroup = m.selectedGroup.Next(scrollDelta, m.showSuccessful)
}
