package box

import tea "github.com/charmbracelet/bubbletea"

func (m *model) selectPrevious(hasMetFirst, hasMetLast bool) {
	if hasMetFirst && hasMetLast {
		return
	}
	m.selectedElement--
	if m.selectedElement < 0 {
		hasMetFirst = true
		m.selectedElement = 0
		if m.elements[m.selectedElement].View() == "" {
			m.selectNext(hasMetFirst, hasMetLast)
		}
	} else if m.elements[m.selectedElement].View() == "" {
		m.selectPrevious(hasMetFirst, hasMetLast)
	}
}

func (m *model) selectNext(hasMetFirst, hasMetLast bool) {
	if hasMetFirst && hasMetLast {
		return
	}
	m.selectedElement++
	if m.selectedElement >= len(m.elements) {
		hasMetLast = true
		m.selectedElement = len(m.elements) - 1
		if m.elements[m.selectedElement].View() == "" {
			m.selectPrevious(hasMetFirst, hasMetLast)
		}
	} else if m.elements[m.selectedElement].View() == "" {
		m.selectNext(hasMetFirst, hasMetLast)
	}
}

func (m *model) cursorFirst() tea.Cmd {
	var cmd1, cmd2 tea.Cmd
	m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
	m.selectedElement = 0
	m.elements[m.selectedElement], cmd2 = m.elements[m.selectedElement].Update(ElementSelectedMsg{true})
	return tea.Batch(cmd1, cmd2)
}

func (m *model) cursorLast() tea.Cmd {
	var cmd1, cmd2 tea.Cmd
	m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
	m.selectedElement = len(m.elements) - 1
	m.elements[m.selectedElement], cmd2 = m.elements[m.selectedElement].Update(ElementSelectedMsg{true})
	return tea.Batch(cmd1, cmd2)
}

func (m *model) cursorWindowUp() tea.Cmd {
	var cmd1, cmd2 tea.Cmd
	m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
	m.selectedElement = m.viewLineToElement[0].idx
	m.elements[m.selectedElement], cmd2 = m.elements[m.selectedElement].Update(ElementSelectedMsg{true})
	return tea.Batch(cmd1, cmd2)
}

func (m *model) cursorWindowDown() tea.Cmd {
	var cmd1, cmd2 tea.Cmd
	m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
	m.selectedElement = m.viewLineToElement[len(m.viewLineToElement)-1].idx
	m.elements[m.selectedElement], cmd2 = m.elements[m.selectedElement].Update(ElementSelectedMsg{true})
	return tea.Batch(cmd1, cmd2)
}

func (m *model) cursorUp() tea.Cmd {
	var cmd1, cmd2 tea.Cmd
	m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
	m.selectPrevious(false, false)
	m.elements[m.selectedElement], cmd2 = m.elements[m.selectedElement].Update(ElementSelectedMsg{true})
	return tea.Batch(cmd1, cmd2)
}
func (m *model) cursorDown() tea.Cmd {
	var cmd1, cmd2 tea.Cmd
	m.elements[m.selectedElement], cmd1 = m.elements[m.selectedElement].Update(ElementSelectedMsg{false})
	m.selectNext(false, false)
	m.elements[m.selectedElement], cmd2 = m.elements[m.selectedElement].Update(ElementSelectedMsg{true})
	return tea.Batch(cmd1, cmd2)
}
