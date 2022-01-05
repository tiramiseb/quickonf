package box

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tiramiseb/quickonf/internal/program/common/style"
	"github.com/tiramiseb/quickonf/internal/program/specific/separator"
)

func (m *model) windowSize(msg tea.WindowSizeMsg) tea.Cmd {
	newMsg := tea.WindowSizeMsg{
		Height: msg.Height - 1,
		Width:  msg.Width - 2,
	}
	m.width = msg.Width
	m.boxHeight = newMsg.Height - 2
	if m.boxHeight < 0 {
		m.boxHeight = 0
	}
	m.updateActive()
	cmds := make([]tea.Cmd, len(m.elements))
	for i, g := range m.elements {
		m.elements[i], cmds[i] = g.Update(newMsg)
	}
	return tea.Batch(cmds...)
}

func (m *model) updateActive() {
	if m.active {
		m.subtitleStyle = style.ActiveBoxTitle.Copy().Width(m.width)
		m.boxStyle = style.ActiveBox.Copy().Width(m.width - 2).Height(m.boxHeight)
	} else {
		m.subtitleStyle = style.BoxTitle.Copy().Width(m.width)
		m.boxStyle = style.Box.Copy().Width(m.width - 2).Height(m.boxHeight)
	}
}

// redrawContent draws content for all elements, even if they are not displayed
func (m *model) redrawContent() {
	result := make([]string, 0, len(m.elements)*2)
	lineToElement := make([]elementLine, 0, len(m.elements)*2)
	for i, e := range m.elements {
		elemView := e.View()
		if elemView == "" {
			continue
		}
		elementView := strings.Split(elemView, "\n")
		thisLineToElement := make([]elementLine, len(elementView))
		for j := range thisLineToElement {
			thisLineToElement[j] = elementLine{i, j}
		}
		result = append(result, elementView...)
		lineToElement = append(lineToElement, thisLineToElement...)
	}
	m.allElementsView = result
	m.allLineToElement = lineToElement
}

// updateView extracts lines to display from content for all elements
func (m *model) updateView() {
	// Find where is the selected element
	startOfSelectedElement := -1
	sizeOfSelectedElement := -1
findElement:
	for i, elementline := range m.allLineToElement {
		switch {
		case startOfSelectedElement < 0 && elementline.idx == m.selectedElement:
			startOfSelectedElement = i
		case startOfSelectedElement >= 0 && elementline.idx != m.selectedElement:
			sizeOfSelectedElement = i - startOfSelectedElement
			break findElement
		}
	}

	// Put the selected element at the top of the view
	firstLineInView := startOfSelectedElement
	lastLineInView := firstLineInView + m.boxHeight

	// Slide up if at the bottom of the list
	if lastLineInView >= len(m.allLineToElement) {
		firstLineInView -= (lastLineInView - len(m.allLineToElement))
		if firstLineInView < 0 {
			firstLineInView = 0
		}
		lastLineInView = firstLineInView + m.boxHeight
		if lastLineInView >= len(m.allLineToElement) {
			lastLineInView = len(m.allLineToElement)
		}
	}

	// Check if the view could be slided up to put the element in the middle
	idealLinesAbove := (m.boxHeight - sizeOfSelectedElement) / 2
	idealFirstLineInVew := startOfSelectedElement - idealLinesAbove
	if idealFirstLineInVew < 0 {
		idealFirstLineInVew = 0
	}
	if idealFirstLineInVew < firstLineInView {
		delta := firstLineInView - idealFirstLineInVew
		firstLineInView = idealFirstLineInVew
		lastLineInView -= delta
	}

	// Extract the view
	m.view = strings.Join(m.allElementsView[firstLineInView:lastLineInView], "\n")
	m.viewLineToElement = m.allLineToElement[firstLineInView:lastLineInView]
	m.selectedElementFirstLine = startOfSelectedElement - firstLineInView
}

func (m *model) cursorPosition() tea.Msg {
	return separator.CursorMsg{
		PointingRight: m.cursorPointsRight,
		Position:      m.selectedElementFirstLine,
	}
}

func (m *model) View() string {
	if len(m.elements) == 0 {
		return m.subtitleStyle.Render(m.title) + "\n" + m.boxStyle.Render(m.msgIfEmpty)
	}
	return m.subtitleStyle.Render(m.title) + "\n" + m.boxStyle.Render(m.view)
}
